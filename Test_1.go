// Licensed Materials - Property of IBM
// (C) Copyright IBM Corp. 2021. All Rights Reserved.
// US Government Users Restricted Rights - Use, duplication or disclosure
// restricted by GSA ADP Schedule Contract with IBM Corp.

package migrations

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mongo_mocks "github.kyndryl.net/MCMP-CommonServices/common-lib-go/pkg/db/mongo/mocks"
	cssettings "github.kyndryl.net/MCMP-CommonServices/common-lib-go/pkg/settings"
	"github.kyndryl.net/MCMP-CommonServices/common-reporting-service/settings"
)
func TestReportIndexes(t *testing.T){
	t.parallel()
	tests := []struct{
		name:			string
		findError		error
		creationError   error
		finalError		error
		model 			[]mongodriver.IndexModel

	}{{
		name: "normal flow",
	},
	{
		name:	"CreateIndex error",
		creationError:errors.New("error"),
		finalError:errors.New("error"),
	},
	{
		name:       "find error",
		findError:  errors.New("error"),
		finalErr:   errors.New("error"),
		model : mongodriver.IndexModel{
			Keys: bson.M{
				name: 1,
			},
	},}

	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, db := mongo_mocks.MockDB(t, viper.GetString(cssettings.ServiceDBNameKey))
			coll := &mongo_mocks.CollectionIfc{}
			db.On("Collection", settings.ReportsCollection).Return(coll)
			cur := &mongo_mocks.CursorIfc{}
			coll.On("IndexExists", ctx, mock.Anything).Return(cur, tt.findError)
			coll.On("CreateOne", ctx, mock.Anything).Return(nil, tt.creationError)
			err := ReportIndexes(ctx, db)
			if tt.name != "normal flow" {
				require.Equal(t, err, tt.finalErr)
			}
		})
	}		
	
}


func TestRemoveDeletedReports(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		decodeData  interface{}
		deleteError error
		findError   error
		finalErr    error
	}{
		{
			name: "normal flow",
		},
		{
			name:        "deleteMany error",
			deleteError: errors.New("error"),
			finalErr:    errors.New("error"),
		},
		{
			name:       "find error",
			decodeData: nil,
			findError:  errors.New("error"),
			finalErr:   errors.New("error"),
		},
		{
			name:       "decode error",
			decodeData: errors.New("error"),
			finalErr:   errors.New("error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, db := mongo_mocks.MockDB(t, viper.GetString(cssettings.ServiceDBNameKey))
			coll := &mongo_mocks.CollectionIfc{}
			db.On("Collection", settings.ReportsCollection).Return(coll)
			cur := &mongo_mocks.CursorIfc{}
			coll.On("Find", ctx, mock.Anything, mock.Anything).Return(cur, tt.findError)
			cur.On("Next", ctx).Return(true).Once()
			cur.On("Next", ctx).Return(false).Once()
			cur.On("Decode", mock.Anything).Return(tt.decodeData)
			coll.On("DeleteMany", ctx, mock.Anything).Return(nil, tt.deleteError)
			err := RemoveDeletedReports(ctx, db)
			if tt.name != "normal flow" {
				require.Equal(t, err, tt.finalErr)
			}
		})
	}
}

func TestAddMetadataToReports(t *testing.T) {
	// type args struct {
	// 	ctx context.Context
	// 	db  mongo.DatabaseIfc
	// }
	tests := []struct {
		name        string
		decodeData  interface{}
		ModifyError error
		findError   error
		finalErr    error
	}{
		{
			name:"normal flow",
		},
		{
			name:       "find error",
			decodeData: nil,
			findError:  errors.New("error"),
			finalErr:   errors.New("error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, db := mongo_mocks.MockDB(t, viper.GetString(cssettings.ServiceDBNameKey))
			coll := &mongo_mocks.CollectionIfc{}
			db.On("Collection", settings.ReportsCollection).Return(coll)
			cur := &mongo_mocks.CursorIfc{}
			coll.On("Find", ctx, mock.Anything, mock.Anything).Return(cur, tt.findError)
			cur.On("Next", ctx).Return(true).Once()
			cur.On("Next", ctx).Return(false).Once()
			cur.On("Decode", mock.Anything).Return(tt.decodeData)
			coll.On("BulkWrite", ctx, mock.Anything).Return(nil, tt.ModifyError)
			err := AddMetadataToReports(ctx, db)
			if tt.name != "normal flow" {
				require.Equal(t, err, tt.finalErr)
			}
		})
	}
}



func TestCopyReportsToS3(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		decodeData  interface{}
		copyError error
		findError   error
		finalErr    error
		DecodeErr	error
		//expected	mongo.ErrSkipMigration
	}{
		{
			name: "normal flow",
		},
		{
			name:        "copyobj error",
			copyError: errors.New("error"),
			finalErr:    errors.New("error"),
		},
		{
			name:		"decodedata error",
			DecodeErr:	errors.New("error"),
			finalErr:	errors.New("error"),

		},
		
		{
			name:       "find error",
			decodeData: nil,
			findError:  errors.New("error"),
			finalErr:   errors.New("error"),
		},
		{
			name:       "decode error",
			decodeData: errors.New("error"),
			finalErr:   errors.New("error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, db := mongo_mocks.MockDB(t, viper.GetString(cssettings.ServiceDBNameKey))
			coll := &mongo_mocks.CollectionIfc{}
			//coll.On("ctx",settings.IsS3Enabled).Return(nil,tt.Expected)
			db.On("Collection", settings.ReportsCollection).Return(coll)
			cur := &mongo_mocks.CursorIfc{}
			coll.On("Find", ctx, mock.Anything, mock.Anything).Return(cur, tt.findError)
			cur.On("Next", ctx).Return(true).Once()
			cur.On("Next", ctx).Return(false).Once()
			cur.On("Decode", mock.Anything).Return(tt.decodeData)
			cur.On("DecodeString", mock.Anything).Return(cur,tt.DecodeErr)
			coll.On("UploadFileToS3", ctx, mock.Anything).Return(nil, tt.copyError)
			err := CopyReportsToS3(ctx, db)
			if tt.name != "normal flow" {
				require.Equal(t, err, tt.finalErr)
			}
		})
	}
}

