package handlers

import (
	"context"
	"net/http"

	"github.com/manujelko/building-go-microservices/pkg/data"
	currencypb "github.com/manujelko/go-grpc-microservices/pkg/protobufs/currency"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//   200: productsResponse

// GetProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()

	// get exchange rate
	req := currencypb.GetRateRequest{
		Base:        currencypb.Currencies(currencypb.Currencies_value["EUR"]),
		Destination: currencypb.Currencies(currencypb.Currencies_value["GPB"]),
	}

	res, err := p.cc.GetRate(context.Background(), &req)
	if err != nil {
		msg := "[Error] error getting new rate"
		p.l.Println(msg, err)
		http.Error(rw, msg, http.StatusInternalServerError)
	}

	for _, p := range lp {
		p.Price = p.Price * res.GetRate()
	}

	err = lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
