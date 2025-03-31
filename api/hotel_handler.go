package api

import (
	"github.com/0x0Glitch/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type HotelHandler struct{
	store 	*db.Store
}

func NewHotelHandler(store *db.Store)*HotelHandler{
	return &HotelHandler{
		store: store,
	}
}

// type HotelQueryParams struct{
// 	Rooms bool
// 	Rating int
// }
func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error{
	id := c.Params("id")
	oid,err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return err
	}
	filter := bson.M{"hotelID": oid}
	rooms,err := h.store.Room.GetRooms(c.Context(),filter)
	if err != nil{
		return err
	}
	return c.JSON(rooms)
	


	return nil
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error{
	
	// if err := c.QueryParser(&qparams); err!= nil{
	// 	return err
	// }
	// fmt.Println(qparams)
	hotels, err := h.store.Hotel.GetHotels(c.Context(),nil)
	if err != nil{
		return err
	}
	return c.JSON(hotels)
}


func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error{
	id := c.Params("id")
	oid ,err := primitive.ObjectIDFromHex(id)
	if err != nil{
		return err
	}
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(),oid)
	return c.JSON(hotel)
}
