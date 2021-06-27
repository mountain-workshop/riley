package bot

import (
	"context"
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun/dbfixture"
)

var testEnv = flag.String("env", "test", "the environment to use when testing")

func StartTestApp(t *testing.T) (context.Context, *App) {
	ctx, app, err := Start(context.TODO(), "test", *testEnv, "", "")
	require.NoError(t, err)
	return ctx, app
}

func loadFixture(t *testing.T, app *App) *dbfixture.Fixture {
	db := app.DB()
	db.RegisterModel((*Team)(nil))

	fixture := dbfixture.New(db, dbfixture.WithTruncateTables())
	err := fixture.Load(app.Context(), FS(), "fixture/test-data.yaml")
	require.NoError(t, err)

	return fixture
}
