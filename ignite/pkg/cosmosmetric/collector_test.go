package cosmosmetric

import (
	"context"
	"errors"
	"testing"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosmetric/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCollector(t *testing.T) {
	// Arrange
	var (
		savedTXs [][]cosmosclient.TX

		fromHeight int64 = 1
	)

	txs := [][]cosmosclient.TX{{}, {}}

	client := mocks.NewTXsCollecter(t)
	client.EXPECT().
		CollectTXs(
			mock.Anything,
			fromHeight,
			mock.AnythingOfType("chan<- []cosmosclient.TX"),
		).
		Run(func(ctx context.Context, fromHeight int64, tc chan<- []cosmosclient.TX) {
			defer close(tc)

			// Send the collected block transactions
			tc <- txs[0]
			tc <- txs[1]
		}).
		Return(nil).
		Times(1)

	db := mocks.NewSaver(t)
	db.EXPECT().
		Save(
			mock.Anything,
			mock.AnythingOfType("[]cosmosclient.TX"),
		).
		Run(func(ctx context.Context, txs []cosmosclient.TX) {
			// Save the transactions
			savedTXs = append(savedTXs, txs)
		}).
		Return(nil).
		Times(2)

	c := NewCollector(db, client)
	ctx := context.Background()

	// Act
	err := c.Collect(ctx, fromHeight)

	// Assert
	require.NoError(t, err)
	require.Equal(t, savedTXs, txs)
}

func TestCollectorWithCollectError(t *testing.T) {
	// Arrange
	wantErr := errors.New("expected error")

	client := mocks.NewTXsCollecter(t)
	client.EXPECT().
		CollectTXs(
			mock.Anything,
			mock.AnythingOfType("int64"),
			mock.AnythingOfType("chan<- []cosmosclient.TX"),
		).
		Run(func(ctx context.Context, fromHeight int64, tc chan<- []cosmosclient.TX) {
			close(tc)
		}).
		Return(wantErr).
		Times(1)

	db := mocks.NewSaver(t)
	c := NewCollector(db, client)
	ctx := context.Background()

	// Act
	err := c.Collect(ctx, 1)

	// Assert
	require.ErrorIs(t, err, wantErr)

	db.AssertNotCalled(t, "Save", mock.Anything, mock.AnythingOfType("[]cosmosclient.TX"))
}

func TestCollectorWithSaveError(t *testing.T) {
	// Arrange
	wantErr := errors.New("expected error")
	txs := []cosmosclient.TX{}

	client := mocks.NewTXsCollecter(t)
	client.EXPECT().
		CollectTXs(
			mock.Anything,
			mock.AnythingOfType("int64"),
			mock.AnythingOfType("chan<- []cosmosclient.TX"),
		).
		Run(func(ctx context.Context, fromHeight int64, tc chan<- []cosmosclient.TX) {
			defer close(tc)

			// Send the collected block transactions
			tc <- txs
		}).
		Return(nil).
		Times(1)

	db := mocks.NewSaver(t)
	db.EXPECT().
		Save(
			mock.Anything,
			mock.AnythingOfType("[]cosmosclient.TX"),
		).
		Return(wantErr).
		Times(1)

	c := NewCollector(db, client)
	ctx := context.Background()

	// Act
	err := c.Collect(ctx, 1)

	// Assert
	require.ErrorIs(t, err, wantErr)
}
