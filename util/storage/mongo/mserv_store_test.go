package mongo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"github.com/TykTechnologies/mserv/storage"
)

var errMsg = "test error"

func Test_getByKey(t *testing.T) {
	t.Parallel()

	eval := is.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("get middleware, with error", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked response.
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11,
					Message: errMsg,
				},
			),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Execute getByKey().
		_, err := m.getByKey(
			context.Background(),
			"key",
			"val",
		)

		// Assert errors.
		eval.Equal(err.Error(), errMsg) // Expected test error.
	})

	mt.Run("get middleware, with success", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked responses.
		doc := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "mw", Value: bson.D{
				{Key: "apiid", Value: "api-1"},
				{Key: "orgid", Value: "org-1"},
				{Key: "uid", Value: "uid-1"},
			}},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.mserv", mtest.FirstBatch, doc),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Execute getByKey().
		res, err := m.getByKey(
			context.Background(),
			"key",
			"val",
		)

		// Assert responses.
		eval.NoErr(err) // Expected no errors.
		eval.Equal(res.APIID, "api-1")
		eval.Equal(res.OrgID, "org-1")
		eval.Equal(res.UID, "uid-1")
	})
}

func Test_GetMWByID(t *testing.T) {
	t.Parallel()

	eval := is.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("get middleware, with error", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked response.
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11,
					Message: errMsg,
				},
			),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Execute GetMWByID().
		_, err := m.GetMWByID(
			context.Background(),
			"id",
		)

		// Assert errors.
		eval.Equal(err.Error(), errMsg) // Expected test error.
	})

	mt.Run("get middleware, with success", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked responses.
		doc := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "mw", Value: bson.D{
				{Key: "apiid", Value: "api-1"},
				{Key: "orgid", Value: "org-1"},
				{Key: "uid", Value: "uid-1"},
			}},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.mserv", mtest.FirstBatch, doc),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Execute GetMWByID().
		res, err := m.GetMWByID(
			context.Background(),
			"id",
		)

		// Assert responses.
		eval.NoErr(err) // Expected no errors.
		eval.Equal(res.APIID, "api-1")
		eval.Equal(res.OrgID, "org-1")
		eval.Equal(res.UID, "uid-1")
	})
}

func Test_GetMWByAPIID(t *testing.T) {
	t.Parallel()

	eval := is.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("get middleware, with error", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked response.
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11,
					Message: errMsg,
				},
			),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Execute GetMWByAPIID().
		_, err := m.GetMWByAPIID(
			context.Background(),
			"id",
		)

		// Assert errors.
		eval.Equal(err.Error(), errMsg) // Expected test error.
	})

	mt.Run("get middleware, with success", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked responses.
		doc := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "mw", Value: bson.D{
				{Key: "apiid", Value: "api-1"},
				{Key: "orgid", Value: "org-1"},
				{Key: "uid", Value: "uid-1"},
			}},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.mserv", mtest.FirstBatch, doc),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Execute GetMWByAPIID().
		res, err := m.GetMWByAPIID(
			context.Background(),
			"id",
		)

		// Assert responses.
		eval.NoErr(err) // Expected no errors.
		eval.Equal(res.APIID, "api-1")
		eval.Equal(res.OrgID, "org-1")
		eval.Equal(res.UID, "uid-1")
	})
}

func Test_GetAllActive(t *testing.T) {
	t.Parallel()

	eval := is.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("get all middlewares, with error", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked response.
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11,
					Message: errMsg,
				},
			),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Execute GetAllActive().
		_, err := m.GetAllActive(context.Background())

		// Assert errors.
		eval.Equal(err.Error(), "collect error: "+errMsg) // Expected test error.
	})

	mt.Run("get all middlewares, with cursor error", func(mt *mtest.T) {
		mt.Parallel()

		// Create invalid subscription to force cursor error.
		doc := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "mw", Value: "invalid"},
		}

		// Defined mocked response.
		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.mserv", mtest.FirstBatch, doc),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Execute GetAllActive().
		_, err := m.GetAllActive(context.Background())

		// Assert response.
		eval.Equal(
			err.Error(),
			"cursor fetch error: error decoding key mw: cannot decode string into a storage.MW",
		) // Expected decoding error.
	})

	mt.Run("get all middlewares, with success", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked responses.
		doc := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "mw", Value: bson.D{
				{Key: "apiid", Value: "api-1"},
				{Key: "orgid", Value: "org-1"},
				{Key: "uid", Value: "uid-1"},
			}},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.mserv", mtest.FirstBatch, doc),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Execute GetAllActive().
		res, err := m.GetAllActive(context.Background())

		// Assert responses.
		eval.NoErr(err) // Expected no errors.
		eval.Equal(res[0].APIID, "api-1")
		eval.Equal(res[0].OrgID, "org-1")
		eval.Equal(res[0].UID, "uid-1")
	})
}

func Test_CreateMW(t *testing.T) {
	t.Parallel()

	eval := is.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("create middleware, with empty id", func(mt *mtest.T) {
		mt.Parallel()

		// Create test store.
		m := Store{}

		// Execute CreateMW().
		_, err := m.CreateMW(context.Background(), &storage.MW{})

		// Assert errors.
		eval.Equal(err, ErrEmptyUID) // Expected empty UID error.
	})

	mt.Run("create middleware, with error", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked response.
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11,
					Message: errMsg,
				},
			),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		mw := storage.MW{
			UID: uuid.NewString(),
		}

		// Execute CreateMW().
		_, err := m.CreateMW(context.Background(), &mw)

		// Assert errors.
		eval.Equal(err.Error(), "insert error: "+errMsg) // Expected test error.
	})

	mt.Run("create middleware, with success", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked responses.
		doc := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "mw", Value: bson.D{
				{Key: "apiid", Value: "api-1"},
				{Key: "orgid", Value: "org-1"},
				{Key: "uid", Value: "uid-1"},
			}},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.mserv", mtest.FirstBatch, doc),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		mw := storage.MW{
			UID: "uid-1",
		}

		// Execute CreateMW().
		id, err := m.CreateMW(context.Background(), &mw)

		// Assert responses.
		eval.NoErr(err) // Expected no errors.
		eval.Equal(id, "uid-1")
	})
}

func Test_UpdateMW(t *testing.T) {
	t.Parallel()

	eval := is.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("update middleware, error empty id", func(mt *mtest.T) {
		mt.Parallel()

		// Create test store.
		m := Store{}

		// Execute UpdateMW().
		_, err := m.UpdateMW(context.Background(), &storage.MW{})

		// Assert errors.
		eval.Equal(err, ErrEmptyUID) // Expected empty UID error.
	})

	mt.Run("update middleware, with find error", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked response.
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11,
					Message: errMsg,
				},
			),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Test middleware.
		mw := storage.MW{
			UID: "uid-1",
		}

		// Execute UpdateMW().
		_, err := m.UpdateMW(context.Background(), &mw)

		// Assert errors.
		eval.Equal(err.Error(), "find error: "+errMsg) // Expected test error.
	})

	mt.Run("update middleware, with update error", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked responses.
		doc := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "mw", Value: bson.D{
				{Key: "apiid", Value: "api-1"},
				{Key: "orgid", Value: "org-1"},
				{Key: "uid", Value: "uid-1"},
			}},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.mserv", mtest.FirstBatch, doc),
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11,
					Message: errMsg,
				},
			),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Test middleware.
		mw := storage.MW{
			UID: "uid-1",
		}

		// Execute UpdateMW().
		_, err := m.UpdateMW(context.Background(), &mw)

		// Assert errors.
		eval.Equal(err.Error(), "update error: "+errMsg) // Expected test error.
	})

	mt.Run("update middleware, with error not found", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked responses.
		doc := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "mw", Value: bson.D{
				{Key: "apiid", Value: "api-1"},
				{Key: "orgid", Value: "org-1"},
				{Key: "uid", Value: "uid-1"},
			}},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.mserv", mtest.FirstBatch, doc),
			mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 0}),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Test middleware.
		mw := storage.MW{
			UID: "uid-1",
		}

		// Execute UpdateMW().
		_, err := m.UpdateMW(context.Background(), &mw)

		// Assert errors.
		eval.Equal(err, ErrNotFound) // Expected error not found.
	})

	mt.Run("update middleware, with success", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked responses.
		doc := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "mw", Value: bson.D{
				{Key: "apiid", Value: "api-1"},
				{Key: "orgid", Value: "org-1"},
				{Key: "uid", Value: "uid-1"},
			}},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.mserv", mtest.FirstBatch, doc),
			mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 1}),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Test middleware.
		mw := storage.MW{
			UID: "uid-1",
		}

		// Execute UpdateMW().
		id, err := m.UpdateMW(context.Background(), &mw)

		// Assert errors.
		eval.NoErr(err) // Expected no error.
		eval.Equal(id, "uid-1")
	})
}

func Test_DeleteMW(t *testing.T) {
	t.Parallel()

	eval := is.New(t)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("delete middleware, error empty id", func(mt *mtest.T) {
		mt.Parallel()

		// Create test store.
		m := Store{}

		// Execute DeleteMW().
		err := m.DeleteMW(context.Background(), "")

		// Assert errors.
		eval.Equal(err, ErrEmptyUID) // Expected empty UID error.
	})

	mt.Run("delete middleware, with find error", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked response.
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11,
					Message: errMsg,
				},
			),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Execute DeleteMW().
		err := m.DeleteMW(context.Background(), "uid-1")

		// Assert errors.
		eval.Equal(err.Error(), "find error: "+errMsg) // Expected test error.
	})

	mt.Run("delete middleware, with delete error", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked responses.
		doc := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "mw", Value: bson.D{
				{Key: "apiid", Value: "api-1"},
				{Key: "orgid", Value: "org-1"},
				{Key: "uid", Value: "uid-1"},
			}},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.mserv", mtest.FirstBatch, doc),
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11,
					Message: errMsg,
				},
			),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Execute DeleteMW().
		err := m.DeleteMW(context.Background(), "uid-1")

		// Assert errors.
		eval.Equal(err.Error(), "delete error: "+errMsg) // Expected test error.
	})

	mt.Run("delete middleware, with success", func(mt *mtest.T) {
		mt.Parallel()

		// Define mocked responses.
		doc := bson.D{
			{Key: "_id", Value: primitive.NewObjectID()},
			{Key: "mw", Value: bson.D{
				{Key: "apiid", Value: "api-1"},
				{Key: "orgid", Value: "org-1"},
				{Key: "uid", Value: "uid-1"},
			}},
		}

		mt.AddMockResponses(
			mtest.CreateCursorResponse(0, "test.mserv", mtest.FirstBatch, doc),
			mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 1}),
		)

		// Create test store.
		m := Store{
			db: mt.Client.Database("test"),
		}

		// Execute DeleteMW().
		err := m.DeleteMW(context.Background(), "uid-1")

		// Assert errors.
		eval.NoErr(err) // Expected no error.
	})
}
