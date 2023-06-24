package query

type Schema struct {
	Name        string        `json:"name,omitempty"`
	Columns     []*Column     `json:"columns,omitempty"`
	ForeignKeys []*ForeignKey `json:"foreign_keys,omitempty"`
}

type Column struct {
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	Length        int64  `json:"length,omitempty"`
	Primary       bool   `json:"primary,omitempty"`
	Unique        bool   `json:"unique,omitempty"`
	NotNull       bool   `json:"not_null,omitempty"`
	Precision     int64  `json:"precision,omitempty"`
	Scale         int64  `json:"scale,omitempty"`
	Default       string `json:"default,omitempty"`
	AutoIncrement bool   `json:"auto_increment,omitempty"`
	Raw           string `json:"raw,omitempty"`
}

type ForeignKey struct {
	Column        string `json:"column,omitempty"`
	RefTable      string `json:"ref_table,omitempty"`
	RefColumn     string `json:"ref_column,omitempty"`
	OnDelCascade  bool   `json:"on_del_cascade,omitempty"`
	OnDelSetNull  bool   `json:"on_del_set_null,omitempty"`
	OnDelSetValue string `json:"on_del_set_value,omitempty"`
}
