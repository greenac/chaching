package helpers

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/greenac/chaching/internal/database/mocks"
	genErr "github.com/greenac/chaching/internal/error"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCreateTable(t *testing.T) {
	Convey("TestCreateTable", t, func() {
		Convey("TestCreateTable should succeed", func() {
			err := CreateTable(context.Background(), mocks.ClientMock{CreateTableOutput: dynamodb.CreateTableOutput{}})
			So(err, ShouldBeNil)
		})

		Convey("TestCreateTable should fail with client error", func() {
			e := errors.New("oak and walnut")
			err := CreateTable(context.Background(), mocks.ClientMock{CreateTableError: e})
			So(err, ShouldResemble, &genErr.GenError{Messages: []string{"CreateTable:Failed to create table with error: " + e.Error()}})
		})
	})
}

func TestDeleteTable(t *testing.T) {
	Convey("TestDeleteTable", t, func() {
		Convey("TestDeleteTable should succeed", func() {
			err := DeleteTable(context.Background(), mocks.ClientMock{DeleteTableOutput: dynamodb.DeleteTableOutput{}}, "table")
			So(err, ShouldBeNil)
		})

		Convey("TestDeleteTable should fail with client error", func() {
			e := errors.New("oak and walnut")
			err := DeleteTable(context.Background(), mocks.ClientMock{DeleteTableError: e}, "table")
			So(err, ShouldResemble, &genErr.GenError{Messages: []string{"DeleteTable:Failed to delete table with error: " + e.Error()}})
		})
	})
}
