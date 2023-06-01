/**
 * Copyright (c) [2023] [Jim-shop]
 * [AircraftWarBackend] is licensed under Mulan PubL v2.
 * You can use this software according to the terms and conditions of the Mulan PubL v2.
 * You may obtain a copy of Mulan PubL v2 at:
 *          http://license.coscl.org.cn/MulanPubL-2.0
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PubL v2 for more details.
 */

package main

import (
	"imshit/aircraftwar/daemon"
	"imshit/aircraftwar/db"
	"imshit/aircraftwar/utils"
	"log"
)

func main() {
	utils.LoadConfig()
	db.InitRedis()
	db.InitSql()
	go daemon.GetPairingDaemon().Run()
	go daemon.GetFightingDaemon().Run()
	r := SetupServer()
	if err := r.RunTLS(":443", "key.pem", "key.key"); err != nil {
		log.Printf("Server run error: %v\n", err)
		return
	}
}
