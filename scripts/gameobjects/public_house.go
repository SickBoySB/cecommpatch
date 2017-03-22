gameobject "public_house" inherit "office"
<<
	local 
	<<
		
		function updateBoozeCounts()
			local warning_data = {}
			
			send("rendUIManager", "SetOfficeInt", SELF, "brewQuantity", #state.brewTable)
			send("rendUIManager", "SetOfficeInt", SELF, "spiritsQuantity", #state.spiritsTable)
			send("rendUIManager", "SetOfficeInt", SELF, "laudanumQuantity", #state.laudanumTable)
			
			if state.distribution_mode == "brew" then
				if #state.brewTable == 0 then
					send("rendUIManager", "SetOfficeString", SELF, "noBoozeWarning", "")
					SELF.tags.has_contents = nil
					table.insert(warning_data,"Brewed Booze")
				else
					send("rendUIManager", "SetOfficeString", SELF, "noBoozeWarning", "")
					SELF.tags.has_contents = true
				end
			elseif state.distribution_mode == "spirits" then
				if #state.spiritsTable == 0 then
					send("rendUIManager", "SetOfficeString", SELF, "noBoozeWarning", "")
					SELF.tags.has_contents = nil
					table.insert(warning_data,"Distilled Spirits")
				else
					send("rendUIManager", "SetOfficeString", SELF, "noBoozeWarning", "")
					SELF.tags.has_contents = true
				end
			end
			
			if state.laudanum_distribution == true then
				if #state.laudanumTable == 0 then
					send("rendUIManager", "SetOfficeString", SELF, "noLaudanumWarning", "")
					SELF.tags.has_laudanum = nil
					table.insert(warning_data,"Laudanum")
				else
					send("rendUIManager", "SetOfficeString", SELF, "noLaudanumWarning", "")
					SELF.tags.has_laudanum = true
				end
			else
				send("rendUIManager", "SetOfficeString", SELF, "noLaudanumWarning", "")
				SELF.tags.has_laudanum = true
			end
			
			send("rendUIManager", "SetOfficeString", SELF, "workPointsStatus", combined_warning_status(warning_data))
		end
		
		function dumpBooze( entityname )
			local results = query("scriptManager",
						"scriptCreateGameObjectRequest",
						"item",
						{ legacyString = entityname} )[1]
						
			send(results,"ClaimItem")
			
			local positionResult = query(SELF, "GetRandomBuildingPosition")[1]
			local x = positionResult.x
			local y = positionResult.y
		
			-- drop outside foundation so we don't get floaters.
			local isInvalidDrop = true
			local i = 0
			while isInvalidDrop do
				if i > 0 then
					positionResult.x = x + rand(i * -1,i)
					positionResult.y = y + rand(i * -1,i)
				end
				isInvalidDrop = query("gameSpatialDictionary","gridHasSpatialTag",positionResult,"occupiedByStructure" )[1]
				i = i + 1
			end
			send(results,"GameObjectPlace",positionResult.x,positionResult.y  )
		end

		function refreshBoozeCapacity()
			local pubData = EntityDB["Public House"]
			
			state.laudanumCapacity = pubData.laudanumCapacityPerVat * state.vatCount
			state.brewCapacity = pubData.brewCapacityPerVat * state.vatCount
			state.spiritsCapacity = pubData.spiritsCapacityPerVat * state.vatCount

			send("rendUIManager", "SetOfficeInt", SELF, "brewCapacity", state.brewCapacity)
			send("rendUIManager", "SetOfficeInt", SELF, "spiritsCapacity", state.spiritsCapacity)
			send("rendUIManager", "SetOfficeInt", SELF, "laudanumCapacity", state.laudanumCapacity)
			
			if #state.laudanumTable >= state.laudanumCapacity then
				SELF.tags.collect_laudanum = nil
				-- dump extra booze
				local dumpAmount = #state.laudanumTable - state.laudanumCapacity
				if dumpAmount > 0 then
					for i=1, dumpAmount do
						dumpBooze( table.remove(state.laudanumTable) )
					end
				end
			else
				if state.laudanum_distribution then
					SELF.tags.collect_laudanum = true
				end
			end
			
			if #state.brewTable >= state.brewCapacity then
				SELF.tags.collect_brew = nil
				-- dump extra booze
				local dumpAmount = #state.brewTable - state.brewCapacity
				if dumpAmount > 0 then
					for i=1, dumpAmount do
						dumpBooze( table.remove(state.brewTable) )
					end
				end
			else
				if state.distribution_mode == "brew" then
					SELF.tags.collect_brew = true
				end
			end
			
			if #state.spiritsTable >= state.spiritsCapacity then
				SELF.tags.collect_spirits = nil
				-- dump extra booze
				local dumpAmount = #state.spiritsTable - state.spiritsCapacity
				if dumpAmount > 0 then
					for i=1, dumpAmount do
						dumpBooze( table.remove(state.spiritsTable) )
					end
				end
			else
				if state.distribution_mode == "spirits" then
					SELF.tags.collect_spirits = true
				end
			end
			
			send("rendUIManager", "SetOfficeInt", SELF, "brewQuantity", #state.brewTable)
			send("rendUIManager", "SetOfficeInt", SELF, "spiritsQuantity", #state.spiritsTable)
			send("rendUIManager", "SetOfficeInt", SELF, "laudanumQuantity", #state.laudanumTable)
		end
	>>

	state
	<<
	>>

	receive Create( stringstringMapHandle init )
	<<
		state.laudanumTable = {}
		state.brewTable = {}
		state.spiritsTable = {}
		
		state.laudanumCapacity = 0
		state.brewCapacity = 0
		state.spiritsCapacity = 0
		
		state.vatCount = 0
		
		state.laudanum_distribution = true
		state.distribution_mode = "brew"
	>>
	
	receive odinBuildingCompleteMessage( int handle, gameSimJobInstanceHandle ji )
	<<
		send("rendUIManager", "SetOfficeInt", SELF, "vatCount",0)
		
		send("rendUIManager", "SetOfficeInt", SELF, "brewQuantity", #state.brewTable)
		send("rendUIManager", "SetOfficeInt", SELF, "spiritsQuantity", #state.spiritsTable)
		send("rendUIManager", "SetOfficeInt", SELF, "laudanumQuantity", #state.laudanumTable)
		
		send("rendUIManager", "SetOfficeInt", SELF, "brewCapacity",0)
		send("rendUIManager", "SetOfficeInt", SELF, "spiritsCapacity", 0)
		send("rendUIManager", "SetOfficeInt", SELF, "laudanumCapacity", 0)
		
		send("rendUIManager", "SetOfficeString", SELF, "noBoozeVatWarning", "")
		send("rendUIManager", "SetOfficeString", SELF, "noChairsWarning", "")
		send("rendUIManager", "SetOfficeString", SELF, "noWorkcrewWarning", "")
		send("rendUIManager", "SetOfficeString", SELF, "noBoozeWarning", "")
		send("rendUIManager", "SetOfficeString", SELF, "noLaudanumWarning", "")
		
		-- default warning/status text
		send("rendUIManager", "SetOfficeString", SELF, "noChairsWarning", "Booze Vats and Chairs needed!")
		send("rendUIManager", "SetOfficeString", SELF, "workPointsStatus", "Brewed Booze and Laudanum needed!")
		
		send("rendUIManager", "SetOfficeString", SELF, "lastGroupSpecialTreatment", "")
		
		-- set up default distribution modes
		send("rendUIManager", "SetOfficeString", SELF, "boozeMode", "Serving: Brewed Drinks")
		send("rendUIManager", "SetOfficeString", SELF, "laudanumMode", "Administering Laudanum on Demand")
		
		--send("rendUIManager", "SetOfficeString", SELF, "boozeModeTooltip", "Now Stocking and Serving: Brewed Drinks")
		
		state.distribution_mode = "brew"
		state.laudanum_distribution = true
		SELF.tags.collect_brew = true
		SELF.tags.collect_spirits = nil
		SELF.tags.collect_laudanum = true
		SELF.tags.has_contents = nil
		updateBoozeCounts()
		refreshBoozeCapacity()
		
		send(SELF,"InteractiveMessage","mode_button1")
	>>
	
	receive InteractiveMessage( string messagereceived )
	<<
		--printl("buildings", "pub received InteractiveMessage: " .. tostring(messagereceived) )
		if messagereceived == "mode_button1" then
			
			state.distribution_mode = "brew"
			SELF.tags.collect_brew = true
			SELF.tags.collect_spirits = nil
			
			send("rendUIManager", "SetOfficeString", SELF, "boozeMode", "Serving: Brewed Drinks")
			--send("rendUIManager", "SetOfficeString", SELF, "boozeModeTooltip", "Now Stocking and Serving: Brewed Drinks")
			updateBoozeCounts()
			refreshBoozeCapacity()
		elseif messagereceived == "mode_button2" then
			
			state.distribution_mode = "spirits"
			SELF.tags.collect_brew = nil
			SELF.tags.collect_spirits = true
			
			send("rendUIManager", "SetOfficeString", SELF, "boozeMode", "Serving: Distilled Spirits")
			--send("rendUIManager", "SetOfficeString", SELF, "boozeModeTooltip", "Now Serving: Brewed Drinks")
			updateBoozeCounts()
			refreshBoozeCapacity()
		--[[elseif messagereceived == "mode_button3" then
			
			state.distribution_mode = "laudanum"
			SELF.tags.collect_brew = nil
			SELF.tags.collect_spirits = nil
			SELF.tags.collect_laudanum = true
			
			send("rendUIManager", "SetOfficeString", SELF, "boozeMode", "Stocking and Serving: Laudanum")
			--send("rendUIManager", "SetOfficeString", SELF, "boozeModeTooltip", "Now Serving: Brewed Drinks")
			updateBoozeCounts()]]
		elseif messagereceived == "laudanum_mode_button1" then
			
			send("rendUIManager", "SetOfficeString", SELF, "laudanumMode", "Administering Laudanum")
			state.laudanum_distribution = true
			SELF.tags.collect_laudanum = true
			updateBoozeCounts()
			refreshBoozeCapacity()
			
		elseif messagereceived == "laudanum_mode_button2" then
			
			send("rendUIManager", "SetOfficeString", SELF, "laudanumMode", "Not administering Laudanum ")
			state.laudanum_distribution = nil
			SELF.tags.collect_laudanum = nil
			updateBoozeCounts()
			refreshBoozeCapacity()
			
		end
	>>

	receive addBoozeItem( gameObjectHandle g)
	<<
		local tags = query(g,"getTags")[1]
		local name = query(g,"getName")[1]
		
		if tags["container"] == true then			
			local containerCount = query(g, "GetContainerContentsCount")
			if containerCount[1] then
				containerCount = containerCount[1]
			else
				-- container is empty for some reason: destroyed objects inside?
				return "abort"
			end
			
			local name = query( query(g, "GetItemInContainer", 1)[1], "getName")[1]
			local tags = query( query(g, "GetItemInContainer", 1)[1], "getTags")[1]
			
			local containerTable = {}
			for i=1,containerCount do
				local item = query(g, "GetItemInContainer", i)[1]
				containerTable[#containerTable + 1] = item
			end
			
			for i = 1,containerCount do
				if tags.laudanum then
					state.laudanumTable[ #state.laudanumTable + 1] = name
				elseif tags.spirits then
					state.spiritsTable[ #state.spiritsTable + 1] = name
				else
					state.brewTable[ #state.brewTable + 1] = name
				end
			end
		else
		
			if tags.laudanum then
				state.laudanumTable[ #state.laudanumTable + 1] = name
			elseif tags.spirits then
				state.spiritsTable[ #state.spiritsTable + 1] = name
			else
				state.brewTable[ #state.brewTable + 1] = name
			end
		end

		updateBoozeCounts()
		refreshBoozeCapacity()
	>>
	
	respond getPubDrinkType()
	<<
		return "getPubDrinkTypeMessage", state.distribution_mode
	>>
	
	respond removeBoozeItem( string boozetype )
	<<
		local name = false
		if boozetype == "brew" then
			if #state.brewTable > 0 then
				name = table.remove(state.brewTable)
			end
		elseif boozetype == "spirits" then
			if #state.spiritsTable > 0 then
				name = table.remove(state.spiritsTable)
			end
		elseif boozetype == "laudanum" then
			if #state.laudanumTable > 0 then
				name = table.remove(state.laudanumTable)
			end
		end
		
		updateBoozeCounts()
		refreshBoozeCapacity()
		
		return "removeBoozeItemMessage", name
	>>
	
	receive recalculateQuality()
	<<
		local has_chairs = false
		local has_booze_vat = false
		local num_vats = 0
		local warning_data = {}
		
		local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
		for k,v in pairs(modules) do
			local tags = query(k, "getTags")[1]
			if tags.booze_vat then
				has_booze_vat = true
				num_vats = num_vats + 1
			end
			if tags.chair then
				has_chairs = true
			end
		end
		
		if has_booze_vat then
			send("rendUIManager", "SetOfficeString", SELF, "noBoozeVatWarning", "")
			send("rendUIManager", "SetOfficeInt", SELF, "vatCount", num_vats)
			state.vatCount = num_vats
		else
			send("rendUIManager", "SetOfficeString", SELF, "noBoozeVatWarning", "")
			send("rendUIManager", "SetOfficeInt", SELF, "vatCount",0)
			table.insert(warning_data,"Booze Vats")
			state.vatCount = 0
		end
		
		if has_chairs then
			send("rendUIManager", "SetOfficeString", SELF, "noChairsWarning", "")
		else
			send("rendUIManager", "SetOfficeString", SELF, "noChairsWarning", "")
			table.insert(warning_data,"Chairs"
		end
		
		send("rendUIManager", "SetOfficeString", SELF, "noChairsWarning", combined_warning_status(warning_data))
		
		updateBoozeCounts()
		refreshBoozeCapacity()
	>>
	
	receive setBuildingOwner(gameObjectHandle newOwner)
	<<
		updateBoozeCounts()
		refreshBoozeCapacity()
	>>
>>