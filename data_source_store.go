package main

import (
    "fmt"
    "log"
	"time"
    "net/http"
    "net/url"
    "encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceStore() *schema.Resource {
	return &schema.Resource{
		Read:   resourceStoreRead,
		Schema: map[string]*schema.Schema{
			"address_url_object": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"store_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
            "delivery_minutes": &schema.Schema{
                Type: schema.TypeInt,
                Computed: true,
            },
		},
	}
}

func resourceStoreRead(d *schema.ResourceData, m interface{}) error {
	var client = &http.Client{Timeout: 10 * time.Second}
    address_url_obj := make(map[string]string)
    err := json.Unmarshal([]byte(d.Get("address_url_object").(string)), &address_url_obj)
    if err != nil {
        return err
    }
    line1 := url.QueryEscape(address_url_obj["line1"])
    log.Printf("%s\n", line1)
    line2 := url.QueryEscape(address_url_obj["line2"])
    log.Printf("%s\n", line2);
    log.Printf("%s\n", fmt.Sprintf("https://www.dominos.com.au/dynamicstoresearchapi/getlimitedstores/25/%s%%20%s", line1, line2));
    stores, err := getStores(fmt.Sprintf("https://www.dominos.com.au/dynamicstoresearchapi/getlimitedstores/25/%s%%20%s", line1, line2), client)
    fmt.Println(err);
    if err != nil {
        return err
    }
    if len(stores) == 0 {
        return fmt.Errorf("No stores near the address %#v", address_url_obj)
    }
    d.Set("store_id", stores[0].StoreNo)
    d.Set("delivery_minutes", stores[0].DeliveryLeadTime)
    d.SetId("store")
	return nil
}


type StoresResponse struct {
	Data []Store
}

type Store struct {
	StoreNo  int
	DeliveryLeadTime int
}

func getStores(url string, client *http.Client) ([]Store, error) {
    r, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer r.Body.Close()
	resp := StoresResponse{}

    err = json.NewDecoder(r.Body).Decode(&resp)
    if err != nil {
        return nil, err
    }
    return resp.Data, nil
}
