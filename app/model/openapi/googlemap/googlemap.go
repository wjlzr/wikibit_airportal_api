package googlemap

type FindCoordinateByAddressRequest struct {
	Address string `json:"address" binding:"required"`
}
