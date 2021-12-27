package filters

type Filter struct {
	Date   string `json:"date"`   //formato "dd-mm-aaaa"
	Status string `json:"status"` //debe ser “true” o “false” stringificados o string vacio para que la api pueda diferenciar entre un status válido o una no eleccion del filtro (string vacio)
	Title  string `json:"title"`
}
