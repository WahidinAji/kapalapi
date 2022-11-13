package vessel

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// type SampleA struct {
// 	Name string `json:"name"`
// }
// type Sample struct {
// 	Datas   []string `json:"datas"`
// 	SampleA `json:"sample_a"`
// }

func (d *VesselDeps) Create(c *fiber.Ctx) error {

	//#region
	// // name := new(SampleA)
	// // if err := c.BodyParser(name); err != nil {
	// // 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// // }
	// ina := new(Sample)
	// if err := c.BodyParser(ina); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// }
	// // ina.SampleA = append(ina.SampleA, *name)
	// transponder := new(Transponder)
	// if err := c.BodyParser(transponder); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// }
	// // licenses := new(Licenses)
	// // if err := c.BodyParser(licenses); err != nil {
	// // 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// // }
	// // engines := new(Engines)
	// // if err := c.BodyParser(engines); err != nil {
	// // 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// // }
	// // fishingCapacity := new(FishingCapacity)
	// // if err := c.BodyParser(fishingCapacity); err != nil {
	// // 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// // }
	//#endregion

	vessel := new(Vessel)
	if err := c.BodyParser(vessel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// var licenses []Licenses
	// if len(vessel.Licenses) > 0 {
	// 	for i := 0; i < len(vessel.Licenses); i++ {
	// 		// var licenses Licenses
	// 		// licenses.ID
	// 		licenses = append(licenses, vessel.Licenses[i])
	// 	}
	// }
	// fmt.Print(licenses)
	// return c.JSON(licenses)
	// vessel.Transponder = *transponder
	// vessel.FishingCapacity = *fishingCapacity

	file, err := c.FormFile("preferredImage")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error get form file ": err.Error()})
	}

	path := "upload"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	err = c.SaveFile(file, fmt.Sprintf("./upload/%s", file.Filename))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error save file to disk ": err.Error()})
	}

	vessel.PreferredImage = fmt.Sprint("./upload/" + file.Filename)
	old := vessel.PreferredImage
	new := fmt.Sprint("./upload/owner-images/" + file.Filename)
	err = os.Rename(old, new)
	if err != nil {
		log.Printf("failed to remove file")
	}

	// err = os.Chmod("Vessel.json", 0755)
	// if err != nil {
	// 	log.Printf("failed")
	// }
	// new := "./upload/Vessel.json"
	// err = os.Rename("Vessel.json",new)
	// if err != nil {
	// 	log.Printf("failed to remove file")
	// }
	// return c.JSON(ina)
	userKey := c.Get("secret-key")

	res, err := d.CreateService(c.Context(), *vessel, userKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	// userKey, err := d.GetUserKeyRepo(c.Context(), in.Uuid)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	// }
	// res, err := d.CreateRepo(c.Context(), in.Vessel, *userKey)
	// if err != nil {
	// 	fmt.Printf("Error creating data")
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	// }

	return c.Status(fiber.StatusOK).JSON(res)
}

func (d *VesselDeps) FindByUserKey(c *fiber.Ctx) error {

	userKey := c.Get("secret-key")
	id := c.Get("user-key-id")
	res, err := d.FindByUserKeyService(c.Context(), UserKey{Id: id, Uuid: userKey})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (d *VesselDeps) CreateNew(c *fiber.Ctx) error {
	c.Accepts("application/json")    // "application/json"
	c.Accepts("multipart/form-data") // "application/json"
	c.Accepts("image/png")           // ""
	c.Accepts("png")

	file, err := c.FormFile("preferredImage")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error get form file ": err.Error()})
	}

	fmt.Println(file.Size)
	if file.Size > 3000000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file cannot greater than 3084 kb"})
	}

	path := "upload"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	fileName := fmt.Sprint(time.Now().Format("20060102150405")) + file.Filename
	err = c.SaveFile(file, fmt.Sprintf("./upload/%s", fileName))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error save file to disk ": err.Error()})
	}

	vessel := new(Vessel)
	if err := c.BodyParser(vessel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error w": err.Error()})
	}

	userKey := c.Get("secret-key")
	vessel.PreferredImage = fileName
	res, err := d.CreateNewService(c.Context(), *vessel, userKey)
	if err != nil {
		err = os.Remove(fmt.Sprintf("./upload/%s", fileName))
		if err != nil {
			log.Printf("remove file failed: %v", err)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error e": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (d *VesselDeps) GetVessel(c *fiber.Ctx) error {
	userKey := c.Get("secret-key")
	vesselKey := c.Get("vessel-secret-key")
	res, err := d.GetVesselService(c.Context(), vesselKey, userKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error e": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (d *VesselDeps) UploadPreferredImage(c *fiber.Ctx) error {
	// img := new(PreferredImage)
	// if err := c.BodyParser(img); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// }
	vesselImage, err := c.FormFile("vesselImage")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "failed to get vessel image",
		})
	}

	return c.Status(fiber.StatusOK).JSON(vesselImage.Filename)
}
