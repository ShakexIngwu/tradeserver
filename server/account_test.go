package server

import (
	"reflect"
	"testing"

	"github.com/ShakexIngwu/tradeserver/webull"
)

func TestGetOpenOrders(t *testing.T) {
	var emptyOrder []Order
	type args struct {
		accountID string
		client    *webull.ClientMock
	}
	tests := []struct {
		name    string
		args    args
		want    []Order
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				accountID: "test_ok",
				client:    &webull.ClientMock{},
			},
			want: []Order{
				{
					Action:                    "BUY",
					ComboTickerType:           "stock",
					FilledQuantity:            1,
					LmtPrice:                  100.00,
					OrderID:                   12345,
					OrderType:                 "LMT",
					OutsideRegularTradingHour: true,
					RemainQuantity:            9,
					Status:                    "Working",
					Symbol:                    "AMZN",
					TickerId:                  54321,
					TimeInForce:               "GTC",
					TotalQuantity:             10,
				},
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				accountID: "test_invalid_orders",
				client:    &webull.ClientMock{},
			},
			want:    emptyOrder,
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				accountID: "test_invalid_remain_quantity",
				client:    &webull.ClientMock{},
			},
			want:    emptyOrder,
			wantErr: false,
		},
		{
			name: "test4",
			args: args{
				accountID: "test_invalid_total_quantity",
				client:    &webull.ClientMock{},
			},
			want:    emptyOrder,
			wantErr: false,
		},
		{
			name: "test5",
			args: args{
				accountID: "test_invalid_lmt_price",
				client:    &webull.ClientMock{},
			},
			want:    emptyOrder,
			wantErr: false,
		},
		{
			name: "test6",
			args: args{
				accountID: "test_invalid_filled_quantity",
				client:    &webull.ClientMock{},
			},
			want:    emptyOrder,
			wantErr: false,
		},
		{
			name: "test7",
			args: args{
				accountID: "test_server_error",
				client:    &webull.ClientMock{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOpenOrders(tt.args.accountID, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOpenOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOpenOrders() got = %v, want %v", got, tt.want)
			}
		})
	}
}
