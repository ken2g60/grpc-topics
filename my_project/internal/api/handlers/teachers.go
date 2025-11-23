package handlers

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/grpc_tutorials/my_project/internal/models"
	"github.com/grpc_tutorials/my_project/internal/repositories/mongodb"
	"github.com/grpc_tutorials/my_project/pkg/utils"
	mainapi "github.com/grpc_tutorials/my_project/proto/gen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// fmt.Println(pbVal)
// fmt.Println("ModelVal:", modelVal)
// fmt.Println("pbVal field[0]:", pbVal.NumField())
// fmt.Println("pbVal num field[0]:", pbVal.Type().Field(0).Name)
// fmt.Println("pbVal Field FirstName[2]:", pbVal.Type().Field(2).Name)
// fmt.Println("pbVal num field[3]:", pbVal.Type().Field(3).Name)
// fmt.Println("pbVal num field[4]:", pbVal.Type().Field(4).Name)
// fmt.Println("pbVal num field[5]:", pbVal.Type().Field(5).Name)
// fmt.Println("pbVal num field[6]:", pbVal.Type().Field(6).Name)
// fmt.Println("ModelVal FieldByName:", modelVal.FieldByName(pbVal.Type().Field(2).Name))
// modelVal.FieldByName(pbVal.Type().Field(4).Name).Set(pbVal.Field(4))
// fmt.Println("ModelVal FirstName:", modelVal.FieldByName("FirstName"))

func (s *Server) GetTeachers(ctx context.Context, req *mainapi.GetTeacherRequest) (*mainapi.Teachers, error) {
	// filter
	filter, err := buildFilter(req.Teacher, &models.Teacher{})
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}

	// sort
	sortOptions := buildSortOptions(req.GetSortBy())

	// connect mongodb
	client, err := mongodb.CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}
	defer client.Disconnect(ctx)

	collection := client.Database("school").Collection("teachers")
	var cursor *mongo.Cursor
	if len(sortOptions) < 1 {
		cursor, err = collection.Find(ctx, filter)
	} else {
		cursor, err = collection.Find(ctx, filter, options.Find().SetSort(sortOptions))
	}
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}
	defer cursor.Close(ctx)

	teachers, err := mongodb.DecodeTeachers(ctx, cursor, func() *mainapi.Teacher { return &mainapi.Teacher{} }, func() *models.Teacher { return &models.Teacher{} })
	if err != nil {
		return nil, err
	}

	return &mainapi.Teachers{Teachers: teachers}, nil

}

func (s *Server) AddTeachers(ctx context.Context, req *mainapi.Teachers) (*mainapi.Teachers, error) {

	for _, teacher := range req.GetTeachers() {
		if teacher.Id != "" {
			return nil, status.Error(codes.InvalidArgument, "request is incorrect format: non-empty field are not allowed")
		}
	}

	addedTeachers, err := GetTeachersToDB(ctx, req.GetTeachers())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &mainapi.Teachers{Teachers: addedTeachers}, nil

}

func (s *Server) UpdateTeachers(ctx context.Context, req *mainapi.Teachers) (*mainapi.Teachers, error) {

	client, err := mongodb.CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}

	defer client.Disconnect(ctx)

	var updatedTeachers []*mainapi.Teacher
	for _, teacher := range req.Teachers {

		if teacher.Id == "" {
			return nil, utils.ErrorHandler(errors.New("id cannot be found"), "id cannot be blank")
		}
		modelTeacher := mapPbTeacherToModelTeacher(teacher)

		objId, err := primitive.ObjectIDFromHex(teacher.Id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "internal error")
		}

		// convert modelTeacher into bson document
		modelDoc, err := bson.Marshal(modelTeacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "internal server")
		}

		var updateDoc bson.M
		err = bson.Unmarshal(modelDoc, &updateDoc)
		if err != nil {
			return nil, utils.ErrorHandler(err, "internal error")
		}

		// Remove the _id field from the update document
		delete(updateDoc, "_id")

		_, err = client.Database("school").Collection("teachers").UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": updateDoc})
		if err != nil {
			return nil, utils.ErrorHandler(err, fmt.Sprintf("error updating teacher id: %s", teacher.Id))
		}

		updatedTeacher := mapModelTeacherToPb(modelTeacher)

		updatedTeachers = append(updatedTeachers, updatedTeacher)
	}

	return &mainapi.Teachers{Teachers: updatedTeachers}, nil
}

func GetTeachersToDB(ctx context.Context, teachersFromReq []*mainapi.Teacher) ([]*mainapi.Teacher, error) {
	client, err := mongodb.CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer client.Disconnect(ctx)

	newTeachers := make([]*models.Teacher, len(teachersFromReq))
	// the fields are dem
	for i, pbTeacher := range teachersFromReq {
		newTeachers[i] = mapPbTeacherToModelTeacher(pbTeacher)
	}

	var addedTeachers []*mainapi.Teacher
	for _, teacher := range newTeachers {
		result, err := client.Database("school").Collection("teachers").InsertOne(ctx, teacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error inserting data")
		}

		objectId, ok := result.InsertedID.(primitive.ObjectID)
		if ok {
			teacher.Id = objectId.Hex()
		}

		// dev environment
		pbTeacher := mapModelTeacherToPb(teacher)
		addedTeachers = append(addedTeachers, pbTeacher)
	}
	return addedTeachers, nil
}

func mapModelTeacherToPb(teacher *models.Teacher) *mainapi.Teacher {
	pbTeacher := &mainapi.Teacher{}
	modelVal := reflect.ValueOf(*teacher)
	pbVal := reflect.ValueOf(pbTeacher).Elem()

	for i := 0; i < modelVal.NumField(); i++ {
		modelField := modelVal.Field(i)
		modelFieldType := modelVal.Type().Field(i)
		pbField := pbVal.FieldByName(modelFieldType.Name)
		if pbField.IsValid() && pbField.CanSet() {
			pbField.Set(modelField)
		}
	}
	return pbTeacher
}

func mapPbTeacherToModelTeacher(pbTeacher *mainapi.Teacher) *models.Teacher {
	modelTeacher := models.Teacher{}
	pbVal := reflect.ValueOf(pbTeacher).Elem()
	modelVal := reflect.ValueOf(&modelTeacher).Elem()
	for i := 0; i < pbVal.NumField(); i++ {
		pbField := pbVal.Field(i)
		fieldName := pbVal.Type().Field(i).Name
		modelField := modelVal.FieldByName(fieldName)
		if modelField.IsValid() && modelField.CanSet() {
			modelField.Set(pbField)
		}
	}
	return &modelTeacher
}
