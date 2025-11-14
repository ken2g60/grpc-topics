#!/bin/bash

# Define the project root folder name
PROJECT_ROOT="my_project"

# Create the main project root directory
mkdir -p $PROJECT_ROOT

# Create subdirectories within the project root
mkdir -p $PROJECT_ROOT/cmd
mkdir -p $PROJECT_ROOT/internal
mkdir -p $PROJECT_ROOT/pkg
mkdir -p $PROJECT_ROOT/proto

echo "Project root folders created successfully."