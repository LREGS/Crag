package store

import (
	"testing"

	"github.com/lregs/Crag/models"
	"github.com/stretchr/testify/assert"
)

func TestStoreClimb(t *testing.T) {

	MockStore := returnPrePopulatedMockStore(t, false, false)

	//need to test if the return value is correct
	t.Run("Testing Store Climb", func(t *testing.T) {
		climb := returnPayload()
		_, err := MockStore.Stores.ClimbStore.StoreClimb(climb)
		if err != nil {
			t.Fatalf("failed because of error: %s", err)
		}
	})

	t.Run("Testing Empty Climb", func(t *testing.T) {
		climb := models.ClimbPayload{}
		_, err := MockStore.Stores.ClimbStore.StoreClimb(climb)
		if err == nil {
			t.Fatal("stored empty result")
		}
	})

	t.Run("Testing hueco grade", func(t *testing.T) {
		climb := models.ClimbPayload{Name: "tank", Grade: "v3", CragID: 1}
		_, err := MockStore.Stores.ClimbStore.StoreClimb(climb)
		if err == nil {
			t.Fatal("accepted incorrect grade")
		}
	})

}

func TestGetClimbsByCragId(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, true, false)

	t.Run("Testing Get Climb By Crag Id", func(t *testing.T) {
		CragId := 1

		climbs, err := MockStore.Stores.ClimbStore.GetClimbsByCragId(CragId)
		if err != nil {
			t.Fatalf("store failedr: %s", err)
		}
		assert.Equal(t, returnClimb(), climbs[0])
	})

}

func TestGetAllClimbs(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, true, false)

	t.Run("Testing GetAllClimbs", func(t *testing.T) {
		climbs, err := MockStore.Stores.ClimbStore.GetAllClimbs()
		if err != nil {
			t.Fatalf("could not get data from store because of err: %s", err)
		}
		if len(climbs) > 0 != true {
			t.Fatalf("No climbs were returned")
		}
	})
}

func TestGetClimbById(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, true, false)

	t.Run("Testing GetClimbByID", func(t *testing.T) {

		testClimb := returnClimb()

		returnedClimb, err := MockStore.Stores.ClimbStore.GetClimbById(1)
		if err != nil {
			t.Fatalf("Could not Get climb by Id becuase of error: %s", err)
		}

		assert.Equal(t, testClimb, returnedClimb)

	})
}

func TestUpdateClimb(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, true, false)

	t.Run("Testing Update Climb", func(t *testing.T) {
		testClimb := returnClimb()
		testClimb.Grade = "9a"

		updatedClimb, err := MockStore.Stores.ClimbStore.UpdateClimb(testClimb)
		if err != nil {
			t.Fatalf("Could not update climb because of err: %s", err)
		}

		assert.Equal(t, testClimb, updatedClimb)
	})

	t.Run("Testing Empty Climb", func(t *testing.T) {
		_, err := MockStore.Stores.ClimbStore.UpdateClimb(models.Climb{})
		if err == nil {
			t.Fatal("shouldn't accept empty climb")
		}
	})
}

func TestDeleteClimb(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, true, false)

	t.Run("Testing Delete Climb", func(t *testing.T) {

		deletedClimb, err := MockStore.Stores.ClimbStore.DeleteClimb(1)
		if err != nil {
			t.Fatalf("could not delete climb because of this error: %s", err)
		}

		assert.Equal(t, returnClimb(), deletedClimb)

		_, err = MockStore.Stores.ClimbStore.GetClimbById(1)
		if err == nil {
			t.Fatalf("Climb still exists in db")
		}

	})
}

func returnPayload() models.ClimbPayload {
	return models.ClimbPayload{
		Name:   "Harvey Oswald",
		Grade:  "7a+",
		CragID: 1,
	}
}

func returnClimb() models.Climb {
	return models.Climb{
		Id:     1,
		Name:   "Harvey Oswald",
		Grade:  "7a+",
		CragID: 1,
	}
}
