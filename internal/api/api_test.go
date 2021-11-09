package api

import (
	"context"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/ozonmp/est-water-api/internal/mocks"
	"github.com/ozonmp/est-water-api/internal/model"
	pb "github.com/ozonmp/est-water-api/pkg/est-water-api"
)

func setup(t *testing.T) (
	*mocks.MockService,
	context.Context,
	pb.EstWaterApiServiceServer,
) {
	ctrl := gomock.NewController(t)

	service := mocks.NewMockService(ctrl)

	ctx := context.Background()

	api := NewWaterAPI(service)

	return service, ctx, api
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
			service, ctx, api := setup(t)

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

			service.EXPECT().
				CreateWater(gomock.Eq(ctx), gomock.Eq(dummyWater.Name), gomock.Eq(dummyWater.Model), gomock.Eq(dummyWater.Material), gomock.Eq(dummyWater.Manufacturer), gomock.Eq(dummyWater.Speed)).
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
		waterId uint64
		expectedErrorMessagePart string
	}{
		"Validate error": {waterId: 0, expectedErrorMessagePart: "InvalidArgument"},
		"404 error": {waterId: 1, expectedErrorMessagePart: "water not found"},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			service, ctx, api := setup(t)

			service.EXPECT().
				DescribeWater(gomock.Eq(ctx), gomock.Eq(testCase.waterId)).
				DoAndReturn(func(ctx context.Context, waterId uint64) (*model.Water, error) {
					return nil, nil
				}).AnyTimes()

			req := pb.DescribeWaterV1Request{WaterId: testCase.waterId}
			_, err := api.DescribeWaterV1(ctx, &req)

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), testCase.expectedErrorMessagePart)
		})
	}
}

func TestApiRemoveWater(t *testing.T) {
	testCases := map[string]struct {
		waterId uint64
		expectedErrorMessagePart string
	}{
		"Validate error": {waterId: 0, expectedErrorMessagePart: "InvalidArgument"},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			service, ctx, api := setup(t)

			service.EXPECT().
				RemoveWater(gomock.Eq(ctx), gomock.Eq(testCase.waterId)).
				DoAndReturn(func(ctx context.Context, waterId uint64) error {
					return nil
				}).AnyTimes()

			req := pb.RemoveWaterV1Request{WaterId: testCase.waterId}
			_, err := api.RemoveWaterV1(ctx, &req)

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), testCase.expectedErrorMessagePart)
		})
	}
}
