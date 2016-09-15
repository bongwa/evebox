package elasticsearch

const (
	SORT_ASC = "asc"
	SORT_DESC = "desc"
)

// Type alias for JSON building.
type mapping map[string]interface{}

// Type alias for JSON building.
type list []interface{}

// Response to an Elastic Search ping (GET /)
type PingResponse struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	Version     struct {
			    Number string `json:"number"`
		    } `json:"version"`
	Tagline     string `json:"tagline"`
}

// Response object for generic responses.
type ResponseObject struct {
	val interface{}
}

func (o ResponseObject) Get(key string) ResponseObject {
	val := o.val.(map[string]interface{})[key]
	return ResponseObject{val: val}
}

func (o ResponseObject) Value() interface{} {
	return o.val
}

// A generic response.
type Response struct {
	body map[string]interface{}
}

func (r *Response) Get(key string) (ResponseObject) {
	return ResponseObject{val: r.body[key]}
}

func NewResponse(body map[string]interface{}) *Response {
	return &Response{
		body: body,
	}
}

type SearchResponse struct {
	Took     uint64 `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
			 Failed     uint64 `json:"failed"`
			 Successful uint64 `json:"successful"`
			 Total      uint64 `json:"total"`
		 } `json:"_shards"`
	Hits     struct {
			 Hits []map[string]interface{} `json:"hits"`
		 } `json:"hits"`

	// The raw response.
	raw      string
}

func (sr SearchResponse) Raw() string {
	return sr.raw
}