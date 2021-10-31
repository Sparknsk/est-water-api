package api

import (
	"context"
	"strings"
	"testing"

	"github.com/ozonmp/est-water-api/internal/mocks"
	"github.com/ozonmp/est-water-api/internal/model"
	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (
	*mocks.MockRepo,
	context.Context,
	pb.EstWaterApiServiceServer,
) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepo(ctrl)
	ctx := context.Background()

	api := NewWaterAPI(repo)

	return repo, ctx, api
}

func TestApiCreateWater(t *testing.T) {
	testCases := map[string]struct {
		field string
		value interface{}
		expectedErrorMessagePart string
	}{
		"Validate error name empty": {field: "Name", value: "", expectedErrorMessagePart: "InvalidArgument"},
		"Validate error name small": {field: "Name", value: "n", expectedErrorMessagePart: "InvalidArgument"},
		"Validate error name big": {field: "Name", value: strings.Repeat("long-name", 10), expectedErrorMessagePart: "InvalidArgument"},
		"Validate error speed small": {field: "Speed", value: 0, expectedErrorMessagePart: "InvalidArgument"},
		"Validate error speed big": {field: "Speed", value: 1001, expectedErrorMessagePart: "InvalidArgument"},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			repo, ctx, api := setup(t)

			dummyWater := model.Water{
				Name: "name",
				Model: "model",
				Manufacturer: "manufacturer",
				Material: "material",
				Speed: 100,
			}
			switch testCase.field {
			case "Name":
				dummyWater.Name = testCase.value.(string)
			case "Model":
				dummyWater.Model = testCase.value.(string)
			case "Manufacturer":
				dummyWater.Manufacturer = testCase.value.(string)
			case "Material":
				dummyWater.Material = testCase.value.(string)
			case "Speed":
				dummyWater.Speed = uint32(testCase.value.(int))
			}

			repo.EXPECT().
				CreateWater(gomock.Eq(ctx), gomock.Eq(model.Water{})).
				DoAndReturn(func(ctx context.Context, water *model.Water) error {
					return nil
				}).AnyTimes()

			req := pb.CreateWaterV1Request{
				Name: dummyWater.Name,
				Model: dummyWater.Model,
				Manufacturer: dummyWater.Manufacturer,
				Material: dummyWater.Material,
				Speed: dummyWater.Speed,
			}
			_, err := api.CreateWaterV1(ctx, &req)

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), testCase.expectedErrorMessagePart)
		})
	}
}

func TestApiDescribeWater(t *testing.T) {
	testCases := map[string]struct {
		waterID uint64
		expectedErrorMessagePart string
	}{
		"Validate error": {waterID: 0, expectedErrorMessagePart: "InvalidArgument"},
		"404 error": {waterID: 1, expectedErrorMessagePart: "water not found"},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			repo, ctx, api := setup(t)

			repo.EXPECT().
				DescribeWater(gomock.Eq(ctx), gomock.Eq(testCase.waterID)).
				DoAndReturn(func(ctx context.Context, waterID uint64) (*model.Water, error) {
					return nil, nil
				}).AnyTimes()

			req := pb.DescribeWaterV1Request{WaterId: testCase.waterID}
			_, err := api.DescribeWaterV1(ctx, &req)

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), testCase.expectedErrorMessagePart)
		})
	}
}

func TestApiRemoveWater(t *testing.T) {
	testCases := map[string]struct {
		waterID uint64
		expectedErrorMessagePart string
	}{
		"Validate error": {waterID: 0, expectedErrorMessagePart: "InvalidArgument"},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			repo, ctx, api := setup(t)

			repo.EXPECT().
				RemoveWater(gomock.Eq(ctx), gomock.Eq(testCase.waterID)).
				DoAndReturn(func(ctx context.Context, waterID uint64) error {
					return nil
				}).AnyTimes()

			req := pb.RemoveWaterV1Request{WaterId: testCase.waterID}
			_, err := api.RemoveWaterV1(ctx, &req)

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), testCase.expectedErrorMessagePart)
		})
	}
}
