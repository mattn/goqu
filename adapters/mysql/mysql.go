package mysql

import (
    "github.com/doug-martin/goqu"
)

var (
    placeholder_rune    = '?'
    quote_rune          = '`'
    singlq_quote        = '\''
    default_values_frag = []byte("")
    mysql_true          = []byte("1")
    mysql_false         = []byte("0")
    time_format         = "2006-01-02 15:04:05"
    operator_lookup     = map[goqu.BooleanOperation][]byte{
        goqu.EQ_OP:                []byte("="),
        goqu.NEQ_OP:               []byte("!="),
        goqu.GT_OP:                []byte(">"),
        goqu.GTE_OP:               []byte(">="),
        goqu.LT_OP:                []byte("<"),
        goqu.LTE_OP:               []byte("<="),
        goqu.IN_OP:                []byte("IN"),
        goqu.NOT_IN_OP:            []byte("NOT IN"),
        goqu.IS_OP:                []byte("IS"),
        goqu.IS_NOT_OP:            []byte("IS NOT"),
        goqu.LIKE_OP:              []byte("LIKE BINARY"),
        goqu.NOT_LIKE_OP:          []byte("NOT LIKE BINARY"),
        goqu.I_LIKE_OP:            []byte("LIKE"),
        goqu.NOT_I_LIKE_OP:        []byte("NOT LIKE"),
        goqu.REGEXP_LIKE_OP:       []byte("REGEXP BINARY"),
        goqu.REGEXP_NOT_LIKE_OP:   []byte("NOT REGEXP BINARY"),
        goqu.REGEXP_I_LIKE_OP:     []byte("REGEXP"),
        goqu.REGEXP_NOT_I_LIKE_OP: []byte("NOT REGEXP"),
    }
)

type DatasetAdapter struct {
    *goqu.DefaultAdapter
}

func (me *DatasetAdapter) SupportsReturn() bool {
    return false
}

func (me *DatasetAdapter) SupportsLimitOnDelete() bool {
    return true
}

func (me *DatasetAdapter) SupportsLimitOnUpdate() bool {
    return true
}

func (me *DatasetAdapter) SupportsOrderByOnDelete() bool {
    return true
}

func (me *DatasetAdapter) SupportsOrderByOnUpdate() bool {
    return true
}

func (me *DatasetAdapter) LiteralString(buf *goqu.SqlBuilder, s string) error {
    if buf.IsPrepared {
        return me.PlaceHolderSql(buf, s)
    }
    buf.WriteRune(singlq_quote)
    for _, char := range s {
        if char == '\'' { // single quote: ' -> \'
            buf.WriteString("\\'")
        } else if char == '"' { // double quote: " -> \"
            buf.WriteString("\\\"")
        } else if char == '\\' { // slash: \ -> "\\"
            buf.WriteString("\\\\")
        } else if char == '\n' { // control: newline: \n -> "\n"
            buf.WriteString("\\n")
        } else if char == '\r' { // control: return: \r -> "\r"
            buf.WriteString("\\r")
        } else if char == 0 { // control: NUL: 0 -> "\x00"
            buf.WriteString("\\x00")
        } else if char == 0x1a { // control: \x1a -> "\x1a"
            buf.WriteString("\\x1a")
        } else {
            buf.WriteRune(char)
        }
    }
    buf.WriteRune(singlq_quote)
    return nil
}

func newDatasetAdapter(ds *goqu.Dataset) goqu.Adapter {
    def := goqu.NewDefaultAdapter(ds).(*goqu.DefaultAdapter)
    def.PlaceHolderRune = placeholder_rune
    def.IncludePlaceholderNum = false
    def.QuoteRune = quote_rune
    def.DefaultValuesFragment = default_values_frag
    def.True = mysql_true
    def.False = mysql_false
    def.TimeFormat = time_format
    def.BooleanOperatorLookup = operator_lookup
    return &DatasetAdapter{def}
}


func init() {
	goqu.RegisterAdapter("mysql", newDatasetAdapter)
}
