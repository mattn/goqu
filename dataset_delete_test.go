package goqu

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func (me *datasetTest) TestDeleteSqlNoReturning() {
	t := me.T()
	mDb, _ := sqlmock.New()
	ds1 := New("no-return", mDb).From("items")
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	_, err := ds1.Returning("id").DeleteSql()
	assert.EqualError(t, err, "goqu: Adapter does not support RETURNING clause")
}

func (me *datasetTest) TestDeleteSqlWithLimit() {
	t := me.T()
	mDb, _ := sqlmock.New()
	ds1 := New("limit", mDb).From("items")
	sql, err := ds1.Limit(10).DeleteSql()
	assert.Nil(t, err)
	assert.Equal(t, sql, `DELETE FROM "items" LIMIT 10`)
}

func (me *datasetTest) TestDeleteSqlWithOrder() {
	t := me.T()
	mDb, _ := sqlmock.New()
	ds1 := New("order", mDb).From("items")
	sql, err := ds1.Order(I("name").Desc()).DeleteSql()
	assert.Nil(t, err)
	assert.Equal(t, sql, `DELETE FROM "items" ORDER BY "name" DESC`)
}

func (me *datasetTest) TestDeleteSql() {
	t := me.T()
	ds1 := From("items")
	sql, err := ds1.DeleteSql()
	assert.NoError(t, err)
	assert.Equal(t, sql, `DELETE FROM "items"`)
}

func (me *datasetTest) TestDeleteSqlNoSources() {
	t := me.T()
	ds1 := From("items")
	_, err := ds1.From().DeleteSql()
	assert.EqualError(t, err, "goqu: No source found when generating delete sql")
}

func (me *datasetTest) TestDeleteSqlWithWhere() {
	t := me.T()
	ds1 := From("items")
	sql, err := ds1.Where(I("id").IsNotNull()).DeleteSql()
	assert.NoError(t, err)
	assert.Equal(t, sql, `DELETE FROM "items" WHERE ("id" IS NOT NULL)`)
}

func (me *datasetTest) TestDeleteSqlWithReturning() {
	t := me.T()
	ds1 := From("items")
	sql, err := ds1.Returning("id").DeleteSql()
	assert.NoError(t, err)
	assert.Equal(t, sql, `DELETE FROM "items" RETURNING "id"`)

	sql, err = ds1.Returning("id").Where(I("id").IsNotNull()).DeleteSql()
	assert.NoError(t, err)
	assert.Equal(t, sql, `DELETE FROM "items" WHERE ("id" IS NOT NULL) RETURNING "id"`)
}

func (me *datasetTest) TestTruncateSql() {
	t := me.T()
	ds1 := From("items")
	sql, err := ds1.TruncateSql()
	assert.NoError(t, err)
	assert.Equal(t, sql, `TRUNCATE "items"`)
}

func (me *datasetTest) TestTruncateSqlNoSources() {
	t := me.T()
	ds1 := From("items")
	_, err := ds1.From().TruncateSql()
	assert.EqualError(t, err, "goqu: No source found when generating truncate sql")
}

func (me *datasetTest) TestTruncateSqlWithOpts() {
	t := me.T()
	ds1 := From("items")
	sql, err := ds1.TruncateWithOptsSql(TruncateOptions{Cascade: true})
	assert.NoError(t, err)
	assert.Equal(t, sql, `TRUNCATE "items" CASCADE`)

	sql, err = ds1.TruncateWithOptsSql(TruncateOptions{Restrict: true})
	assert.NoError(t, err)
	assert.Equal(t, sql, `TRUNCATE "items" RESTRICT`)

	sql, err = ds1.TruncateWithOptsSql(TruncateOptions{Identity: "restart"})
	assert.NoError(t, err)
	assert.Equal(t, sql, `TRUNCATE "items" RESTART IDENTITY`)

	sql, err = ds1.TruncateWithOptsSql(TruncateOptions{Identity: "continue"})
	assert.NoError(t, err)
	assert.Equal(t, sql, `TRUNCATE "items" CONTINUE IDENTITY`)
}

func (me *datasetTest) TestPreparedDeleteSql() {
	t := me.T()
	ds1 := From("items")
	sql, args, err := ds1.ToDeleteSql(true)
	assert.NoError(t, err)
	assert.Equal(t, args, []interface{}{})
	assert.Equal(t, sql, `DELETE FROM "items"`)

	sql, args, err = ds1.Where(I("id").Eq(1)).ToDeleteSql(true)
	assert.NoError(t, err)
	assert.Equal(t, args, []interface{}{1})
	assert.Equal(t, sql, `DELETE FROM "items" WHERE ("id" = ?)`)

	sql, args, err = ds1.Returning("id").Where(I("id").Eq(1)).ToDeleteSql(true)
	assert.NoError(t, err)
	assert.Equal(t, args, []interface{}{1})
	assert.Equal(t, sql, `DELETE FROM "items" WHERE ("id" = ?) RETURNING "id"`)
}

func (me *datasetTest) TestPreparedTruncateSql() {
	t := me.T()
	ds1 := From("items")
	sql, args, err := ds1.ToTruncateSql(true, TruncateOptions{})
	assert.NoError(t, err)
	assert.Equal(t, args, []interface{}{})
	assert.Equal(t, sql, `TRUNCATE "items"`)
}
