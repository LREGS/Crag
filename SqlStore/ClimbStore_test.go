package store

import (
	"reflect"
	"testing"

	"github.com/lregs/Crag/models"
)

func TestAddClimb(t *testing.T) {

	MockStore := returnPrePopulatedMockStore(t, false, false)

	//need to test if the return value is correct
	t.Run("Testing Add Climb", func(t *testing.T) {
		climb := returnClimb()
		_, err := MockStore.Stores.ClimbStore.StoreClimb(climb)
		if err != nil {
			t.Fatalf("failed because of error: %s", err)
		}
	})

}

func TestGetClimbsByCrag(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, true, false)

	t.Run("Testing Get Climb", func(t *testing.T) {
		CragId := 1

		climbs, err := MockStore.Stores.ClimbStore.GetClimbsByCrag(CragId)
		if err != nil {
			t.Fatalf("could not get data from store because of err: %s", err)
		}
		if climbs[0].Name != "harvey Oswald" {
			t.Fatalf("Returned Climb is named %s but expected Harvey Oswald", climbs[0].Name)
		}
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

		MockClimb := returnClimb()

		climb, err := MockStore.Stores.ClimbStore.GetClimbById(1)
		if err != nil {
			t.Fatalf("Could not Get climb by Id becuase of error: %s", err)
		}

		if reflect.TypeOf(climb) != reflect.TypeOf(MockClimb) {
			t.Fatalf("The returned climb does not equal the standard climb loaded into the db")
		}

	})
}

func TestUpdateClimb(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, true, false)

	t.Run("Testing Update Climb", func(t *testing.T) {
		MockClimb := returnClimb()
		MockClimb.Grade = "9a"

		row, err := MockStore.Stores.ClimbStore.UpdateClimb(MockClimb)
		if err != nil {
			t.Fatalf("Could not update climb because of err: %s", err)
		}

		if row.Grade != MockClimb.Grade {
			t.Fatalf("Update failed, grade is %s but expected %s", row.Grade, MockClimb.Grade)
		}
		// var updatedClimb models.Climb
		// err = row.Scan(&updatedClimb.Id, &updatedClimb.Name, &updatedClimb.Grade, &updatedClimb.CragID)

		// if reflect.TypeOf(MockClimb) != reflect.TypeOf(updatedClimb) {
		// 	t.Fatalf("Expected update grade to %s but %s was returned", MockClimb.Grade, updatedClimb.Grade)
		// }

	})
}

func TestDeleteClimb(t *testing.T) {
	MockStore := returnPrePopulatedMockStore(t, true, false)

	t.Run("Testing Delete Climb", func(t *testing.T) {
		climb := returnClimb()

		err := MockStore.Stores.ClimbStore.DeleteClimb(climb.Id)
		if err != nil {
			t.Fatalf("could not delete climb because of this error: %s", err)
		}

		_, err = MockStore.Stores.ClimbStore.GetClimbById(climb.Id)
		if err == nil {
			t.Fatalf("Climb still exists in db")
		}

	})
}

func returnClimb() *models.Climb {
	climb := &models.Climb{
		Id:     1,
		Name:   "harvey Oswald",
		Grade:  "7a+",
		CragID: 1,
	}
	return climb
}
