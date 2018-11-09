package ipv6

// List fetches a list of IPv6 models according to the specified (optional) `Filters`.
func List(f Filters) (*ListResult, error) {

	return &ListResult{
		Results: []*Address{
			&Address{
				Network:   "2c0f:ff90::/32",
				Latitude:  1,
				Longitude: 38,
			},
		},
	}, nil
}
