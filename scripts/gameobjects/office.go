gameobject "office" inherit "buildings"
<<
	local 
	<<
	>>

	state
	<<
		int rOH
		string buildingName
		string buildingFancyName
		string buildingTitle
		int buildingQuality
		string buildingQualityName
		table squares
		bool completed
		
		int currentJobIndex
		
		table parent
		table jobs
		gameSimAssignmentHandle curConstructionAssignment

		gameSimAssignmentHandle curAssignment

		int assignmentID
		string currentJob
		string currentProduct

		table materials
		bool claimed
		gameObjectHandle buildingOwner
		
		bool resupply
		table supplies
	>>

	receive Create( stringstringMapHandle init )
	<<
		local myName = init["legacyString"]
		state.entityName = myName
          state.buildingTitle = "Office"
		SELF.tags.overseer_active = nil

		if state.buildingName == "Foreign Office" then
			send("rendUIManager", "SetOfficeInt", SELF, "desksPresent", 0)
		elseif state.buildingName == "Barbershop" then 
			send("rendUIManager", "SetOfficeInt", SELF, "chairsPresent", 0)
		elseif state.buildingName == "Chapel" then
			send("rendUIManager", "SetOfficeInt", SELF, "chairsPresent", 0)
		end
		
		if EntityDB[state.entityName].supply_info then
			state.resupply = true
		else
			state.resupply = false
		end
		
		state.supplies = {}
		local supply_info = EntityDB[state.entityName].supply_info
		if supply_info then
			for k,v in pairs(supply_info) do
				state.supplies[ v.tier ] = 0
			end
		end

		if state.entityName ~= "Trade Office" then
			send("gameSession", "incSessionInt", "workplaceCount", 1)
		end
	>>

	receive setBuildingOwner(gameObjectHandle newOwner)
	<<
		-- w/ new owner, if any modules should push jobs to office, push 'em now.

		printl("buildings", "new building owner set on office, attempting to push automatic module jobs")
		
		if newOwner == nil then
			send("rendUIManager", "SetOfficeBool", SELF, "building_assigned", false)
			SELF.tags.overseer_active = nil
		else
			send("rendUIManager", "SetOfficeBool", SELF, "building_assigned", true)
			local fn = query(newOwner, "getFirstName")[1]
			local ln = query(newOwner, "getLastName")[1]

			send("rendUIManager", "SetOfficeString", SELF, "ownerFirstName", fn)
			send("rendUIManager", "SetOfficeString", SELF, "ownerLastName", ln)
			
			if query(newOwner,"isOnShift")[1] then
				SELF.tags.overseer_active = true
			end
		end
		
		if newOwner then
			send("rendUIManager", "SetOfficeString", SELF, "noOverseerWarning", "")
		else
			send("rendUIManager", "SetOfficeString", SELF, "noOverseerWarning", "")
		end
	>>

	receive DeregisterSupervisor ( gameObjectHandle supervisor )
	<<
		if buildingOwner == supervisor then
			state.claimed = false
			state.buildingOwner = nil
		end
	>>

	receive InteractiveMessage( string messagereceived )
	<<
		printl ("buildings", "Message Received: " .. messagereceived )
		if not state.completed then
			return
		end

		-- OTHER STUFF HERE
		if messagereceived == "build" then
			
			-- Open menu
			if state.slatedForDemolition == true then
				send("rendUIManager", "OpenOfficeMenu", SELF, "Demolished" )
			else
				send("rendUIManager", "OpenOfficeMenu", SELF, state.buildingName )
			end
			
			if state.buildingOwner == nil then
				send("rendUIManager", "SetOfficeMenuHeader", SELF, state.buildingName, "overseer_icon");
			else
				local name = query(state.buildingOwner, "GetWorkPartyName");
				if name then
					send("rendUIManager", "SetOfficeMenuHeader", SELF, state.buildingName, "overseer_icon")
				else
					send("rendUIManager", "SetOfficeMenuHeader", SELF, state.buildingName, "overseer_icon")
				end
			end
		elseif messagereceived == "modules" then
			printl("buildings", "opening modules menu...")
			send("rendCommandManager", "OpenModulesMenu", SELF)
		end
	>>

	receive odinBuildingCompleteMessage ( int handle, gameSimJobInstanceHandle ji )
	<<
		if state.supplies and state.supplies[1] then
			for k,v in pairs(state.supplies) do
				send("rendUIManager", "SetOfficeInt", SELF, "workPoints" .. tostring(k), state.supplies[k])
				--printl("DAVID", " setting office int for tier " .. k )
			end
		end
		
		send("gameWorkshopManager", "AddOffice", SELF, state.buildingName)
		printl("buildings", "office " .. tostring(SELF.id) .. " Got building complete.")
		
		if query("gameSession", "getSessionBool", "chapelBuilt")[1] == false then
			if state.buildingName == "Chapel" then
				send("gameSession", "setSessionBool", "chapelBuilt", true)
			end
		end

		local officeInfo = EntityDB[state.entityName]
		
		if officeInfo and officeInfo.standing_jobs then
			for i=1,#officeInfo.standing_jobs do
				send("gameWorkshopManager", "AddOfficeStandingOrder", SELF, officeInfo.standing_jobs[i])
			end
		end
		
		state.buildingTitle = state.parent.name

		send(SELF,"SetGoodName")
		-- send("rendUIManager", "SetOfficeString", SELF, "eventName", "THE DARK GRIMOIRE");
		
		local buildingIcon = "offices_category"
		if SELF.tags.barracks then
			buildingIcon = "barracks"
			send("rendUIManager", "SetOfficeString", SELF, "backgroundFilenameLarge", "ui\\workcrewbackgrounds\\bg_barracks")
		elseif SELF.tags.chapel then
			buildingIcon = "chapel"
		elseif SELF.tags.barbershop then
			buildingIcon = "barbershop_image"
		elseif SELF.tags.public_house then
			buildingIcon = "messhall"
		elseif SELF.tags.mine then
			buildingIcon = "mineshaft_image"
		elseif SELF.tags.naturalists_office then
			buildingIcon = "exploration_plus_naturalism"
		elseif SELF.tags.foreign_office then
			buildingIcon = "foreign_office"
		elseif SELF.tags.trade_office then
			buildingIcon = "trade_office"
		elseif SELF.tags.laboratory then
			buildingIcon = "laboratory"
		end
		
		if not SELF.tags.trade_office then	
			send("rendCommandManager",
				"odinRendererStubMessage",
				"ui\\orderIcons.xml",
				"artisan_icon",
				state.buildingName .. " built.", -- header text
				"The ".. state.buildingName .. " is constructed and needs an overseer assigned to it.", -- text description
				"Left-click for details. Right-click to dismiss.", -- action string
				"needOwner", -- alert type (for stacking)
				"ui//eventart//heliograph.png", -- imagename for bg
				"high", -- importance: low / high / critical
				state.rOH, -- object ID
				60 * 1000, -- duration in ms
				0, -- "snooze" time if triggered multiple times in rapid succession
				nil) 
		end

		send("rendCommandManager",
			"odinRendererTickerMessage",
			"A " .. state.buildingName .. " finished construction on day " ..
				query("gameSession","getSessionInt","dayCount")[1] .. ".",
			buildingIcon,
			"ui\\orderIcons.xml")
		
		ready()
		
		state.hitpoints =  #state.squares --10 -- you get your hitpoints!
		state.hitpointsmax = #state.squares -- For display and possibly repair
		send("rendUIManager", "SetOfficeInt", SELF, "buildingHP", state.hitpoints)
		send("rendUIManager", "SetOfficeInt", SELF, "buildingHPMax", state.hitpointsmax)
		send("rendUIManager", "SetOfficeString", SELF, "buildingHPDescription", "Undamaged")	
		send("rendUIManager","SetOfficeString",SELF,"buildingFancyName",state.buildingName )
		send(SELF, "recalculateQuality")
		
		if not SELF.tags.trade_office then
			send("rendUIManager", "SetOfficeString", SELF, "noOverseerWarning", "")
			send("rendUIManager", "SetOfficeInt", SELF, "workPoints1", 0)
			SELF.tags.no_supplies1 = true
		end
	>>

	receive SetGoodName()
	<<
		if state.parent == nil then
			printl("buildings", tostring(SELF.id) .. " office.go STATE PARENT IS NIL")
		end
		SELF.tags["good_name"] = true
		SELF.tags["bad_name"] = nil
		state.buildingFancyName = "The " .. goodAdjectives[rand(1,#goodAdjectives)]:gsub("^%l", string.upper) .. " " .. state.parent.good_names[rand(1,#state.parent.good_names)]
		send("rendUIManager","SetOfficeString",SELF,"buildingFancyName",state.buildingFancyName )
		send("rendOdinBuildingClassHandler", "odinSetBuildingName", SELF, state.buildingFancyName);
	>>

	receive SetBadName()
	<<
		SELF.tags["good_name"] = nil
		SELF.tags["bad_name"] = true
		state.buildingFancyName =  "The " .. badAdjectives[rand(1,#badAdjectives)]:gsub("^%l", string.upper) .. " " .. state.parent.bad_names[rand(1,#state.parent.bad_names)]
		send("rendUIManager","SetOfficeString",SELF,"buildingFancyName",state.buildingFancyName )
		send("rendOdinBuildingClassHandler", "odinSetBuildingName", SELF, state.buildingFancyName);
	>>

	receive ownerDied( string reason )
	<<
		local s = ""
		if reason == "death" then
			s = " needs a new overseer assigned because the previous overseer met an unfortunate end."
		elseif reason == "left" then
			s = " needs a new overseer assigned because the previous overseer left the colony."
		else
			s = " needs a new overseer assigned."
		end
		
		send("rendCommandManager",
			"odinRendererStubMessage",
			"ui\\orderIcons.xml",
			"artisan_icon",
			"Overseer Needed", -- header text
			"The ".. state.buildingName .. s, -- text description
			"Left-click for details. Right-click to dismiss.", -- action string
			"needOwner", -- alert type (for stacking)
			"ui//eventart//heliograph.png", -- imagename for bg
			"high", -- importance: low / high / critical
			state.rOH, -- object ID
			60 * 1000, -- duration in ms
			0, -- "snooze" time if triggered multiple times in rapid succession
			nil) -- gameobjecthandle of director, null if none
	>>
		
	receive ModuleFinishedProduction (gameObjectHandle module)
	<<
	
	>>

	receive Update()
	<<

	>>
	
	receive DestroyBuilding( gameSimJobInstanceHandle ji )
	<<
		if not SELF.tags.trade_office then
			send("gameSession", "incSessionInt", "workplaceCount", -1)
		end

	>>
	
	receive AssignmentCancelledMessage( gameSimAssignmentHandle assignment )
	<<
		if not state.completed then
			send("gameSession", "incSessionInt", "workplaceCount", -1)
		end
	>>

	receive recalculateQuality()
	<<

	>>
>>
