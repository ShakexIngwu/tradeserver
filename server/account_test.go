package server

import (
	"reflect"
	"testing"

	"github.com/ShakexIngwu/tradeserver/webull"
)

func TestGetOpenOrders(t *testing.T) {
	var emptyOrder []*order
	type args struct {
		accountID string
		client    *webull.ClientMock
	}
	tests := []struct {
		name    string
		args    args
		want    []*order
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				accountID: "test_ok",
				client:    &webull.ClientMock{},
			},
			want: []*order{
				{
					action:                    "BUY",
					ComboTickerType:           "stock",
					filledQuantity:            1,
					lmtPrice:                  100.00,
					orderID:                   12345,
					orderType:                 "LMT",
					outsideRegularTradingHour: true,
					remainQuantity:            9,
					status:                    "Working",
					symbol:                    "AMZN",
					tickerId:                  54321,
					timeInForce:               "GTC",
					totalQuantity:             10,
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
