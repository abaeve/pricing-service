package model

import "time"

type OrderPayload struct {
	FetchRequestId string `json:"request_id,omitempty" bson:"fetch_request_id"` /* request_id string */
	RegionId       int32  `json:"region_id,omitempty" bson:"region_id"`         /* region_id integer */

	OrderId      int64     `json:"order_id,omitempty" bson:"order_id"`           /* order_id integer */
	TypeId       int32     `json:"type_id,omitempty" bson:"type_id"`             /* type_id integer */
	LocationId   int64     `json:"location_id,omitempty" bson:"location_id"`     /* location_id integer */
	SystemId     int32     `json:"system_id,omitempty" bson:"system_id"`         /* The solar system this order was placed */
	VolumeTotal  int32     `json:"volume_total,omitempty" bson:"volume_total"`   /* volume_total integer */
	VolumeRemain int32     `json:"volume_remain,omitempty" bson:"volume_remain"` /* volume_remain integer */
	MinVolume    int32     `json:"min_volume,omitempty" bson:"min_volume"`       /* min_volume integer */
	Price        float64   `json:"price,omitempty" bson:"price"`                 /* price number */
	IsBuyOrder   bool      `json:"is_buy_order,omitempty" bson:"is_buy_order"`   /* is_buy_order boolean */
	Duration     int32     `json:"duration,omitempty" bson:"duration"`           /* duration integer */
	Issued       time.Time `json:"issued,omitempty" bson:"issued"`               /* issued string */
	Range_       string    `json:"range,omitempty" bson:"range"`                 /* range string */
}
