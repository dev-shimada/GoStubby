package usecase_test

import (
	"reflect"
	"testing"

	"github.com/dev-shimada/gostubby/internal/domain/repository"
	"github.com/dev-shimada/gostubby/internal/usecase"
	"github.com/google/go-cmp/cmp"
)

func TestEndpointUsecase_EndpointMatcher(t *testing.T) {
	type fields struct {
		cr repository.ConfigRepository
	}
	type args struct {
		arg usecase.EndpointMatcherArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.EndpointMatcherResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eu := usecase.NewEndpointUsecase(tt.fields.cr)
			got, err := eu.EndpointMatcher(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("EndpointUsecase.EndpointMatcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EndpointUsecase.EndpointMatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpointUsecase_ResponseCreator(t *testing.T) {
	type fields struct {
		cr repository.ConfigRepository
	}
	type args struct {
		arg usecase.ResponseCreatorArgs
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    usecase.ResponseCreatorResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eu := usecase.NewEndpointUsecase(tt.fields.cr)
			got, err := eu.ResponseCreator(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("EndpointUsecase.ResponseCreator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}
