package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Student struct {
	ID    primitive.ObjectID
	Name  string
	Kelas Kelas
}

type School struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string
}

type Kelas struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Sekolah primitive.ObjectID
	Name    string
}

func NewSchool(name string) *School {
	return &School{
		ID:   primitive.NewObjectID(),
		Name: name,
	}
}

func NewKelas(sekolah primitive.ObjectID, name string) *Kelas {
	return &Kelas{
		ID:      primitive.NewObjectID(),
		Sekolah: sekolah,
		Name:    name,
	}
}
