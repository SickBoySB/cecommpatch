gameobject "barracks" inherit "office"
<<
	local 
	<<
		function barracks_reset_supply_text()
			-- CECOMMPATCH - premature "No Door" alert fix
			if not state.completed then
				return
			end
			
			local supply_warning = ""
			local status = ""
			local weapon_data = EntityDB[ state.weaponLoadout ]
			
			if not state.ammo_stocking_goals then
				send(SELF,"recalculateQuality")
			end
			
			if not state.buildingOwner or not weapon_data.ammo_tier then
				-- no nothing. Owner alert handled automatically.
				-- turn off all jobs.
				for i=1,4 do
					SELF.tags["needs_resupply" .. tostring(i)] = nil
					SELF.tags["needs_resupply" .. tostring(i) .. "_badly"] = nil
					SELF.tags["no_supplies" ..  tostring(i)] = true
				end
			else
				-- only stock what ammo we are using.
				-- only warn for what ammo we are using.
				local ammo = weapon_data.ammo_tier
				
				-- there are 4 types of ammo numbered 1-4. Turn on the one we're using (if any). Turn off everything else.
				for i=1,4 do
					if i == ammo then
						if state.supplies[i] >= state.ammo_stocking_goals[ "ammo" .. tostring(i) ] then
							status = "Working. Stocked with supplies."
							if state.resupply == false then
								status = "Working. Resupply HALTED."
							end
							SELF.tags["needs_resupply" .. tostring(i)] = nil
							SELF.tags["needs_resupply" .. tostring(i) .. "_badly"] = nil
							SELF.tags["no_supplies" ..  tostring(i)] = nil

						elseif state.supplies[i] >= EntityDB[ state.entityName ].mc_resupply_when_below then
							status = "Working. Low supplies."
							if state.resupply == false then
								status = "Working. Resupply HALTED."
								SELF.tags["needs_resupply" .. tostring(i)] = nil
								SELF.tags["needs_resupply" .. tostring(i) .. "_badly"] = nil
							else
								SELF.tags["needs_resupply" .. tostring(i)] = true
								SELF.tags["needs_resupply" .. tostring(i) .. "_badly"] = nil
							end
							SELF.tags["no_supplies" ..  tostring(i)] = nil
							supply_warning = "Running low on ammo."
							
						elseif state.supplies[i] == 0 then
							status = "Using fallback weapon. Out of supplies."
							SELF.tags["no_supplies" ..  tostring(i)] = true
							if state.resupply == false then
								status = "Using fallback weapon. Resupply HALTED."
								SELF.tags["needs_resupply" .. tostring(i)] = nil
								SELF.tags["needs_resupply" .. tostring(i) .. "_badly"] = nil
							else
								SELF.tags["needs_resupply" .. tostring(i)] = true
								SELF.tags["needs_resupply" .. tostring(i) .. "_badly"] = true
								if state.buildingOwner then
									local ownername = query(state.buildingOwner,"getName")[1]
									local alertstring = "The Barracks operated by " .. ownername .. " is out of ammunition! Procure more ammo so soldiers can fight using better weapons!"
												
									send("rendCommandManager",
										"odinRendererStubMessage", --"odinRendererStubMessage",
										"ui\\orderIcons.xml", -- iconskin
										"barracks", -- icon
										"Barracks needs ammo!", -- header text
										alertstring, -- text description
										"Left-click to zoom. Right-click to dismiss.", -- action string
										"barracksProblem", -- alert typasde (for stacking)
										"ui\\eventart\\cult_ritual.png", -- imagename for bg
										"low", -- importance: low / high / critical
										state.rOH, -- object ID
										60 * 1000, -- duration in ms
										0, -- snooze
										state.director)
								end
							end
							supply_warning = "Ammo required."
						end
					else
						SELF.tags["needs_resupply" .. tostring(i)] = nil
						SELF.tags["needs_resupply" .. tostring(i) .. "_badly"] = nil
						SELF.tags["no_supplies" ..  tostring(i)] = true
					end
				end
			end
			
			if not state.buildingOwner then
				status = "Work HALTED. Overseer needed."
			end
			
			send("rendUIManager","SetOfficeString",SELF,"noSuppliesWarning",supply_warning)
			send("rendUIManager","SetOfficeString",SELF,"workPointsStatus",status)
		end
	>>

	state
	<<
		int currentShiftSelection
		string weaponLoadout
	>>

	receive Create( stringstringMapHandle init )
	<<
		state.currentShiftSelection = 1 -- 1 (day) or 5 (evening)
		send(SELF,"setWeaponLoadout","pistol") -- pistol.
		
		send("rendUIManager", "SetOfficeBool", SELF, "musketEnabled", false)
		send("rendUIManager", "SetOfficeBool", SELF, "blunderbussEnabled", false)
		send("rendUIManager", "SetOfficeBool", SELF, "jezailEnabled", false)
		send("rendUIManager", "SetOfficeBool", SELF, "tripistolEnabled", false)
		send("rendUIManager", "SetOfficeBool", SELF, "revolverEnabled", false)
		send("rendUIManager", "SetOfficeBool", SELF, "carbineEnabled", false)
		send("rendUIManager", "SetOfficeBool", SELF, "grenadelauncherEnabled", false)
		send("rendUIManager", "SetOfficeBool", SELF, "leydenpistolEnabled", false)
		send("rendUIManager", "SetOfficeBool", SELF, "leydenrifleEnabled", false)
		
		state.ammo_stocking_goals = {
			ammo1 = 0,
			ammo2 = 0,
			ammo3 = 0,
			ammo4 = 0,
		}
	>>
	
	receive InteractiveMessage( string messagereceived )
	<<
		printl("buildings", "barracks received message: " .. messagereceived)
		if not state.completed then
			return
		end
		
		if messagereceived == "supply_on" then
			state.resupply = true
			barracks_reset_supply_text()
		elseif messagereceived == "supply_off" then
			state.resupply = false
			barracks_reset_supply_text()
		elseif messagereceived == "DayOrNight_button1" then
			state.currentShiftSelection = 1 -- morning
			printl("buildings", "barracks selected day shift")
			
			-- if overseer update overseer's shift
			if state.buildingOwner then
				printl("buildings", "got buildingowner: " .. query(state.buildingOwner,"getName")[1] .. " and setting to day shift" )
				--send(state.buildingOwner,"InteractiveMessage","setStartHour1")
				--send("rendCommandManager", "gameSetWorkPartyWorkShift", state.buildingOwner, 1)
				send(state.buildingOwner,"setWorkShift",  0, 0, 0, 0, 2, 5, 1, 1) 
			end
			
		elseif messagereceived == "DayOrNight_button2" then
			state.currentShiftSelection = 5 -- evening
			printl("buildings", "barracks selected night shift")
			
			-- if overseer update overseer's shift
			if state.buildingOwner then
				printl("buildings", "got buildingowner: " .. query(state.buildingOwner,"getName")[1] .. " and setting to night shift" )
				--send(state.buildingOwner,"InteractiveMessage","setStartHour5")
				--send("rendCommandManager", "gameSetWorkPartyWorkShift", state.buildingOwner, 5)
				send(state.buildingOwner,"setWorkShift",  2, 5, 1, 1, 0, 0, 0, 0)
			end
			
		else
			-- if string has "Loadout_" in it, extract the reset and see if the remainder matches a programName
			local result = string.find(messagereceived, "Loadout_")
			
			if result then
				local loadoutName = string.lower( string.sub( messagereceived, 9 ) )
				
				printl("buildings", " barracks got loadout request for: " .. loadoutName )
				
				-- loadout button string TO entity name (for when they're different)
				local button_to_entity = {
					["jezail"] = "jezail_rifle",
					["grenade launcher"] = "grenadelauncher",
					["leyden pistol"] = "leyden_pistol",
					["leyden rifle"] = "leyden_rifle",
				}
				
				if button_to_entity[loadoutName] then
					loadoutName = button_to_entity[loadoutName]
				end
				
				send(SELF,"setWeaponLoadout", loadoutName)
			end
		end
	>>
	
	receive setWeaponLoadout(string loadoutname)
	<<
		printl("buildings", "barracks got setWeaponLoadout: " .. tostring(loadoutname) )
		if loadoutname then
			state.weaponLoadout = loadoutname
		else
			if not state.weaponLoadout then
				state.weaponLoadout = "pistol"
			end
		end
		
		--[[local ammo_types = {
			"Stone Pellet Ammunition",
			"Ball Cartridge Ammunition",
			"Full Metal Jacket Ammunition",
			"Crate of Grenades",
			"Leyden Jars",
			"Pressurized Petroleum",
		}]]
		
		local weapon_data = EntityDB[state.weaponLoadout]
		
		send("rendUIManager","SetOfficeString",SELF,"weaponName",  tostring( weapon_data.display_name ) )
		
		--send("rendUIManager","SetOfficeString",SELF,"weaponInfo1", tostring( weapon_data.ranged_damage ))
		--send("rendUIManager","SetOfficeString",SELF,"weaponInfo2", tostring( weapon_data.reload_time * 0.1 ) .. " seconds")
		--send("rendUIManager","SetOfficeString",SELF,"weaponInfo3", tostring( ammo_types[ weapon_data.ammo_tier ] ))
		
		barracks_reset_supply_text()
	>>
	
	respond getWeaponLoadout()
	<<
		return "weaponLoadoutMessage", state.weaponLoadout
	>>
	
	receive setBuildingOwner(gameObjectHandle newOwner)
	<<
		-- set new owner to correct shift.
		if newOwner then
			if not state.currentShiftSelection then state.currentShiftSelection = 1 end
			--send(newOwner,"InteractiveMessage","setStartHour" .. tostring(state.currentShiftSelection))
			--send("rendCommandManager", "gameSetWorkPartyWorkShift", newOwner, state.currentShiftSelection)
			
			if state.currentShiftSelection == 1 then
				send(newOwner,"setWorkShift", 0, 0, 0, 0, 2, 5, 1, 1) 
			elseif state.currentShiftSelection == 5 then
				send(newOwner,"setWorkShift", 2, 5, 1, 1, 0, 0, 0, 0)
			end
			send(SELF,"setWeaponLoadout",nil)
		end
	>>
	
	respond GetActiveShift()
	<<
		return "barracksShiftResult", state.currentShiftSelection
	>>
	
	receive recalculateQuality()
	<<
		state.ammo_stocking_goals = {
			ammo1 = 0,
			ammo2 = 0,
			ammo3 = 0,
			ammo4 = 0,
		}
	
		-- do weapon locker detection here; every time a module is added or removed, this is run.
		-- oops entity name isn't even used.
		local tag_to_weapon = {
			musket_locker = {boolName = "musketEnabled", entity="musket",enabled=false, ammo="ammo1"},
			carbine_locker = {boolName = "carbineEnabled", entity="carbine",enabled=false, ammo="ammo3"},
			blunderbuss_locker = {boolName = "blunderbussEnabled", entity="blunderbuss",enabled=false, ammo="ammo2"},
			jezail_locker = {boolName = "jezailEnabled", entity="jezail_rifle",enabled=false, ammo="ammo3"},
			tripistol_locker = {boolName = "tripistolEnabled", entity="tripistol",enabled=false, ammo="ammo1"},
			revolver_locker = {boolName = "revolverEnabled", entity="revolver",enabled=false, ammo="ammo2"},
			grenade_launcher_locker = {boolName = "grenadelauncherEnabled", entity="grenadelauncher",enabled=false, ammo="ammo4"},
			leyden_pistol_locker = {boolName = "leydenpistolEnabled", entity="musket",enabled=false},
			leyden_rifle_locker = {boolName = "leydenrifleEnabled", entity="musket",enabled=false},
		}
		
		-- detect various weaponry crates here.
		-- AND detect number of weapon lockers per ammo type
		local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
		for k,v in pairs(modules) do
			for tagname,tagbool in pairs(v.tags) do
				if tag_to_weapon[tagname] then
					printl("buildings", "barracks found locker: " .. tagname)
					tag_to_weapon[tagname].enabled = true
					if tag_to_weapon[tagname].ammo then
						state.ammo_stocking_goals[ tag_to_weapon[tagname].ammo ] = state.ammo_stocking_goals[ tag_to_weapon[tagname].ammo ] +
																		EntityDB[ state.entityName ].lc_resupply_when_below
					end
				end
			end
		end
		
		for k,v in pairs(tag_to_weapon) do
			if v.enabled == false then
				printl("buildings", "did NOT get : " .. v.boolName )
			end
			send("rendUIManager", "SetOfficeBool", SELF, v.boolName, v.enabled)
		end
		
		barracks_reset_supply_text()
	>>
	
	receive shiftToggle(string status)
	<<
		-- ????

		--[[if status == "off" then
			send("rendUIManager", "SetOfficeInt", SELF, "shiftStatus", "OFF")
		elseif status == "on" then
			send("rendUIManager", "SetOfficeInt", SELF, "shiftStatus", "ON")
		else
			printl("CHRIS", "Office received a bad shift toggle!!")
		end]]
	>>
	
	receive odinBuildingCompleteMessage ( int handle, gameSimJobInstanceHandle ji )
	<<
		send("rendUIManager","SetOfficeString",SELF,"weaponName","")
		--send("rendUIManager","SetOfficeString",SELF,"weaponInfo1","")
		--send("rendUIManager","SetOfficeString",SELF,"weaponInfo2","")
		--send("rendUIManager","SetOfficeString",SELF,"weaponInfo3","")
		
		send("rendUIManager","SetOfficeInt",SELF,"workPoints1",0)
		send("rendUIManager","SetOfficeInt",SELF,"workPoints2",0)
		send("rendUIManager","SetOfficeInt",SELF,"workPoints3",0)
		send("rendUIManager","SetOfficeInt",SELF,"workPoints4",0)
		--send("rendUIManager", "SetOfficeInt", SELF, "workPoints5", 0)
		--send("rendUIManager", "SetOfficeInt", SELF, "workPoints6", 0)
		
		barracks_reset_supply_text()
		
		
		if not query("gameSession","getSessionBool","builtABarracks")[1] then
			send("gameSession","setSessionBool","builtABarracks",true)
			send("gameSession", "setSteamAchievement", "builtABarracks")
		end
		
		send(SELF,"InteractiveMessage","Loadout_Pistol")
	>>

	receive buildingAddSupplies( gameObjectHandle item, int count)
	<<
		local tags = query(item,"getTags")[1]
		local mult = 1
		local tier = 1

		for k,v in pairs( EntityDB[ state.entityName ].supply_info ) do
			if tags[k] then
				mult = v.multiplier
				tier = v.tier
				break
			end
		end

		state.supplies[tier] = state.supplies[tier] + count * mult

		send("rendUIManager", "SetOfficeInt", SELF, "workPoints" .. tostring(tier) , state.supplies[tier])
		barracks_reset_supply_text()
	>>
	
	receive consumeSupplies( int tier, int count )
	<<
		if not tier or not count then
			return
		end
		
		state.supplies[tier] = state.supplies[tier] - count
		
		-- can consume more ammo than barracks has, but this trick only works once per, so NBD.
		if state.supplies[tier] < 0 then state.supplies[tier] = 0 end
		
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints" .. tostring(tier), state.supplies[tier])
		barracks_reset_supply_text()
	>>
>>