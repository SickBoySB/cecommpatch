
gameobject "mine" inherit "office"
<<
	local 
	<<
          function UpdateMineProducts()

               local localProductsParagraph = "NONE"
               for k,v in pairs(state.localProducts) do
                    if k ~= nil and v == true then  
                         local name = EntityDB[k].display_name
                         if localProductsParagraph == "NONE" then
                              localProductsParagraph = name
                         else
                              localProductsParagraph = localProductsParagraph .. ", " .. name
                         end
                    end
               end
               
               if #state.specialProducts > 0 then
                    for k,v in pairs(state.specialProducts) do
					if localProductsParagraph == "NONE" then
						local name = EntityDB[v].display_name
						localProductsParagraph = name
					else
						local name = EntityDB[v].display_name
						localProductsParagraph = localProductsParagraph .. ", " .. name
					end
                    end
               end
               
               printl("buildings", "local mine products are : " .. localProductsParagraph)
               
               send("rendUIManager", "SetOfficeString", SELF, "localProducts", localProductsParagraph)
               
               local checkTable = { {check = false},
                    {check = true, boolname = "shallowMixedEnabled",  disabledString = "Requires a minimum mine depth of 10."},
                    {check = true, boolname = "coreVeinEnabled", 	disabledString = "Requires a minimum mine depth of 50 and the installation of a Mine Dewatering Pump."},
                    {check = true, boolname = "deepVeinEnabled", 	disabledString = "Requires a minimum mine depth of 200 and the installation of a Ventilation Unit."},
                    {check = true, boolname = "farReachesEnabled", 	disabledString = "Requires a minimum mine depth of 500 and the installation of a Steam Distributor."},
                    }
               
               local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
               local mineshaft = false 
               for k,v in pairs(modules) do
                    if v.tags.mineshaft then 
                         mineshaft = true
                         break
                    end
               end
               
               if mineshaft == true then
                    --This just attaches all the outputs in a text list.
                    
                    --NOTES: This is where obfuscation will happen later. hardcoding it for now.
                    local obfuscationTable = {false,false,false,false,true}
                    
                    for i=1,5 do
                         local paragraph = "NONE"
                         if obfuscationTable[i] == false then
                              --First put in the default products.
                              for k,v in pairs(state.defaultProducts["strata" .. i]) do
                                   local name = EntityDB[v].display_name
                                   if paragraph == "NONE" then
                                        paragraph = name
                                   else
                                       paragraph = paragraph .. ", " .. name
                                   end
                              end
                              
                              --Then local products.
                              if (i > 1) then --add local products
                                   for k,v in pairs(state.localProducts) do
                                        if v == true then
                                             local name = EntityDB[k].display_name
                                             if paragraph == "NONE" then
                                                  paragraph = name
                                             else
                                                  paragraph = paragraph .. ", " .. name
                                             end
                                        end
                                   end
                              end
                              
                              --Special products, if you have them unlocked.
                              if i > 3 and #state.specialProducts <= 0 then
                                   paragraph = paragraph .. ", ???, ???"
                              elseif i > 3 and #state.specialProducts == 1 then
                                   for k,v in pairs(state.specialProducts) do
								local name = EntityDB[v].display_name
								if paragraph == "NONE" then
									paragraph = name
								else
									paragraph = paragraph .. ", " .. name
								end
                                   end
                                   paragraph = paragraph .. ", ???"
                              elseif (i > 3) then
                                   for k,v in pairs(state.specialProducts) do
								local name = EntityDB[v].display_name
								if paragraph == "NONE" then
									paragraph = name
								else
									paragraph = paragraph .. ", " .. name
								end
                                   end
                              end
                         else
                              paragraph = "What can be found here is a mystery." --You get this if obfuscation is true.
                         end
                         send("rendUIManager", "SetOfficeString", SELF, "strata" .. i .. "products", paragraph)
                         --Now let's add a thing detailing whether it's locked or not.
                         local paragraph2 = ""
                         if checkTable[i].check == true then
                              if state[checkTable[i].boolname] ~= true then --Add some text telling you what you need.
                                   paragraph2 = checkTable[i].disabledString
                              end
                         end
                         send("rendUIManager", "SetOfficeString", SELF, "strata" .. i .. "enabledString", paragraph2)
                    end
               else
                    --You don't have a mineshaft. You get no info.
                    for i=1,5 do
                         local paragraph = "Build a mineshaft to see what will be available from this mine."
                         send("rendUIManager", "SetOfficeString", SELF, "strata" .. i .. "products", paragraph)
                         --Now let's add a thing detailing whether it's locked or not.
                         local paragraph2 = "You cannot mine without a mineshaft."
                         send("rendUIManager", "SetOfficeString", SELF, "strata" .. i .. "enabledString", paragraph2)
                    end
               end
          end
          
          function checkForButtonEnablement()
               local checkTable = {
                    {	boolname = "shallowMixedEnabled", checkDepth = 10, 	strataname = "Shallow Mixed", needsModule = false},
                    {	boolname = "coreVeinEnabled", checkDepth = 50, 		strataname = "Core Vein", 	needsModule = true, moduleRequired = "dewateringPump"},
                    {	boolname = "deepVeinEnabled", checkDepth = 200, 		strataname = "Deep Vein", 	needsModule = true, moduleRequired = "ventilationUnit"},
                    {	boolname = "farReachesEnabled", checkDepth = 500, 	strataname = "Far Reaches", 	needsModule = true, moduleRequired = "steamDistributor"},
                    }
               
               for k,v in pairs(checkTable) do
                    --printl("CHRIS", "check is " .. v.boolname .. " and " .. tostring(state[v.boolname]))
                    if (state[v.boolname] == false) and state.mineDepth >= v.checkDepth then
                         if v.needsModule == true then
                              --printl("CHRIS", "module is " .. tostring(v.needsModule) .. " and " .. v.moduleRequired .. " and " .. tostring(state[v.moduleRequired]))
                              if state[v.moduleRequired] == true then
                                   state[v.boolname] = true
                                   send("rendUIManager", "SetOfficeBool", SELF, v.boolname, true)
							
                                   send("rendCommandManager",
                                        "odinRendererStubMessage",
                                        "ui\\orderIcons.xml",
                                        "mineshaft_image",
                                        "Stratum Unlocked!", -- header text
                                        "Our miners have gained access to the " .. v.strataname .. " stratum.", -- text description
                                        "Right-click to dismiss.", -- action string
                                        "miningAlert", -- alert type (for stacking)
                                        "ui//eventart//mineShoring.png", -- imagename for bg
                                        "low", -- importance: low / high / critical
                                        state.rOH, -- object ID
                                        30000, -- duration in ms
                                        30, -- "snooze" time if triggered multiple times in rapid succession
                                        nil) -- gameobjecthandle of director, null if none
                                   
                                   UpdateMineProducts()
                              end
                         else
                              state[v.boolname] = true
                              send("rendUIManager", "SetOfficeBool", SELF, v.boolname, true)
						
                              send("rendCommandManager",
                                   "odinRendererStubMessage",
                                   "ui\\orderIcons.xml",
                                   "mineshaft_image",
                                   "Stratum Unlocked!", -- header text
                                   "Our miners have gained access to the " .. v.strataname .. " stratum.", -- text description
                                   "Right-click to dismiss.", -- action string
                                   "miningAlert", -- alert type (for stacking)
                                   "ui//eventart//mineShoring.png", -- imagename for bg
                                   "low", -- importance: low / high / critical
                                   state.rOH, -- object ID
                                   30000, -- duration in ms
                                   30, -- "snooze" time if triggered multiple times in rapid succession
                                   nil) -- gameobjecthandle of director, null if none
						
                              UpdateMineProducts()
                         end
					
                    elseif (state[v.boolname] == true) and v.needsModule == true then --Let's see if you've REMOVED any modules you need.
                         if state[v.moduleRequired] ~= true then
                              --disable it!!!
                              state[v.boolname] = false
                              send("rendUIManager", "SetOfficeBool", SELF, v.boolname, false)
						
                              send("rendCommandManager",
                                   "odinRendererStubMessage",
                                   "ui\\orderIcons.xml",
                                   "mineshaft_image",
                                   "Stratum lost!", -- header text
                                   "Our miners have lost access to the " .. v.strataname .. " stratum.", -- text description
                                   "Right-click to dismiss.", -- action string
                                   "miningAlert", -- alert type (for stacking)
                                   "ui//eventart//mineShoring.png", -- imagename for bg
                                   "low", -- importance: low / high / critical
                                   state.rOH, -- object ID
                                   30000, -- duration in ms
                                   30, -- "snooze" time if triggered multiple times in rapid succession
                                   nil) -- gameobjecthandle of director, null if none
						
                              UpdateMineProducts()
                         end
                    end
               end
          end
		
		function mine_reset_supply_text()
			
			local status = ""
			local supply_warning = ""
			
			if state.supplies[1] >= EntityDB[ state.entityName ].lc_resupply_when_below then
				status = "Working. Mine is supplied."
				if state.resupply == false then
					status = "Working. Workcrew ordered to NOT re-supply mine."
				end
				SELF.tags.no_supplies1 = nil
				SELF.tags.needs_resupply1 = nil
				SELF.tags.needs_resupply1_badly = nil
				
			elseif state.supplies[1] >= EntityDB[ state.entityName ].mc_resupply_when_below then
				status = "Working. Mine Shorings low. Assigned labourers will re-supply mine."
				if state.resupply == false then
					status = "Working. Low on Mine Shorings. Workcrew ordered to NOT re-supply mine."
					SELF.tags.needs_resupply1 = nil
					SELF.tags.needs_resupply1_badly = nil
				else
					SELF.tags.needs_resupply1 = true
					SELF.tags.needs_resupply1_badly = nil
				end
				SELF.tags.no_supplies1 = nil
				supply_warning = "Running low on Mine Shorings."
				
			elseif state.supplies[1] == 0 then
				status = "Work halted. Mine Shorings needed. Workcrew will re-supply mine."
				if state.resupply == false then
					status = "Work halted. Workcrew ordered to NOT re-supply mine."
					SELF.tags.needs_resupply1 = nil
					SELF.tags.needs_resupply1_badly = nil
				else
					SELF.tags.needs_resupply1 = true
					SELF.tags.needs_resupply1_badly = true

					if state.buildingOwner then

						local ownername = query(state.buildingOwner,"getName")[1]
						local alertstring = "The Mine operated by " .. ownername .. " is out of Mine Shorings! Produce more Mine Shorings so mining can continue."
						
						send("rendCommandManager",
							"odinRendererStubMessage", --"odinRendererStubMessage",
							"ui\\commodityIcons.xml", -- iconskin
							"mine_shoring", -- icon
							"Mine Needs Shorings", -- header text
							alertstring, -- text description
							"Left-click to zoom. Right-click to dismiss.", -- action string
							"mineProblem", -- alert typasde (for stacking)
							"ui\\eventart\\minerCrew.png", -- imagename for bg
							"low", -- importance: low / high / critical
							state.rOH, -- object ID
							60 * 1000, -- duration in ms
							0, -- snooze
							state.director)
					end
				end
				
				SELF.tags.no_supplies1 = true
				supply_warning = "Mine Shorings required to do work."
			end
			
			send("rendUIManager", "SetOfficeString", SELF, "noSuppliesWarning",supply_warning)
			send("rendUIManager", "SetOfficeString", SELF, "workPointsStatus", status)
		end
	>>

	state
	<<
		int mineDepth
	>>

	receive Create( stringstringMapHandle init )
	<<
		state.mineDepth = 1
          
          state.defaultProducts = {}
          state.defaultProducts.strata1 = {"rough_stone_block","cube_of_clay","bushel_of_sand"}
          state.defaultProducts.strata2 = {"rough_stone_block","cube_of_clay","bushel_of_sand","sulphur"}
          state.defaultProducts.strata3 = {"rough_stone_block","coal",}
          state.defaultProducts.strata4 = {"coal"}
          state.defaultProducts.strata5 = {"coal"} --Maybe remove coal here later
		
		--IMPORTANT NOTE: In defaultProducts.strata, the TABLE VALUE is the name of the product. In localProducts, we use the name of the product as the TABLE KEY instead to prevent duplication.
		state.localProducts = {}
		
          --Special ores added post-mining
		state.specialProducts = {}
		
		--This is POTENTIAL special products
          state.lowQualityProductTable = {"sulphur","coal","chalk"}
          state.highQualityProductTable = {"hematite","malachite","sphalerite","native_gold"}
		-- this is for super rare neat/silly veins. They should pop a special notification.
          state.SSRQualityProductTable ={"logs","steel_ingots","obeliskian_block"} --Note to self: add something that lets these enable events somehow so they can be dangerous
		
		-- This is for artifacts and other one-offs. This is never added to the description paragraph.
          state.rareItems = {"bushel_of_bones","bushel_of_scrap_iron","leaf_fossil","serpent_bell","obeliskian_block","boxed_rare_painting","steel_ingots"} 
          
          state.ventilationUnit = false
          state.dewateringPump = false
          state.steamDistributor = false
          
          state.shallowMixedEnabled = false
          state.coreVeinEnabled = false
          state.deepVeinEnabled = false
          state.farReachesEnabled = false
		
          send("rendUIManager", "SetOfficeBool", SELF, "shallowMixedEnabled", false)
		send("rendUIManager", "SetOfficeBool", SELF, "coreVeinEnabled", false)
		send("rendUIManager", "SetOfficeBool", SELF, "deepVeinEnabled", false)
		send("rendUIManager", "SetOfficeBool", SELF, "farReachesEnabled", false)
          
          state.currentStrata = 1
	>>

	receive odinBuildingCompleteMessage ( int handle, gameSimJobInstanceHandle ji )
	<<
		send("rendUIManager","SetOfficeString",SELF,"noShaftWarning","At least one mineshaft is required to perform work.")
		send("rendUIManager", "SetOfficeString", SELF, "mineDepth", tostring(state.mineDepth) .. " yards")
		
		UpdateMineProducts()
		
		send(SELF,"refreshTechModifierDisplay")
		
		mine_reset_supply_text()
	>>
     
     receive InteractiveMessage( string messagereceived )
	<<
		printl("buildings", "mine " .. tostring(SELF.id) .. " got InteractiveMessage: " .. tostring(messagereceived) )
		if not state.completed then
			return
		end
		if messagereceived == "MineStrata_button1" then
			state.currentStrata = 0
		elseif messagereceived == "MineStrata_button2" then
			state.currentStrata = 1
		elseif messagereceived == "MineStrata_button3" then
			state.currentStrata = 2
		elseif messagereceived == "MineStrata_button4" then
			state.currentStrata = 3
          elseif messagereceived == "MineStrata_button5" then
			state.currentStrata = 4
          elseif messagereceived == "MineStrata_button6" then
			state.currentStrata = 5
			
		elseif messagereceived == "supply_on" then
			state.resupply = true
			mine_reset_supply_text()
		elseif messagereceived == "supply_off" then
			state.resupply = false
			mine_reset_supply_text()
		end
	>>
	
     receive recalculateQuality()
	<<
          --reset all module states to prepare for checking
          state.ventilationUnit = false
          state.dewateringPump = false
          state.steamDistributor = false
          local has_shaft = false
		
          local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
		for k,v in pairs(modules) do 
               local moduleName = query(v, "getModuleName")[1]
               if moduleName == "Dewatering Pump" then
                    state.dewateringPump = true
               end
			
               if moduleName == "Ventilation Unit" then
                    state.ventilationUnit = true
               end
			
               if moduleName == "Steam Distributor" then
                    state.steamDistributor = true
               end
			
			if moduleName == "Mineshaft Module" then
				has_shaft = true
			end
		end
          
		if not has_shaft then
			send("rendUIManager","SetOfficeString",SELF,"noShaftWarning","At least one mineshaft is required to perform work.")
		else
			send("rendUIManager","SetOfficeString",SELF,"noShaftWarning","")
		end
			
		send(SELF,"refreshTechModifierDisplay")
          checkForButtonEnablement() --Let's see if we've enabled/disabled any strata.
     >>
	
	respond SupplyMineOutput()
	<<

          local digAmount = 1
          if state.currentStrata == 0 then
               digAmount = 3
          end
          
          if (state.mineDepth + digAmount >= 500) and state.farReachesEnabled ~= true then
               state.mineDepth = 500
               if state.currentStrata == 0 then
                    send("rendCommandManager",
                         "odinRendererStubMessage",
                         "ui\\orderIcons.xml",
                         "mineshaft_image",
                         "Mine Deepening Halted!", -- header text
                         "Your mineshaft has reached 500 feet and requires a Steam Distributor to go any deeper!", -- text description
                         "Right-click to dismiss.", -- action string
                         "miningAlert", -- alert type (for stacking)
                         "ui//eventart//mineShoring.png", -- imagename for bg
                         "low", -- importance: low / high / critical
                         state.rOH, -- object ID
                         30 * 1000, -- duration in ms
                         30, -- "snooze" time if triggered multiple times in rapid succession
                         nil) -- gameobjecthandle of director, null if none
               end
			
          elseif (state.mineDepth + digAmount >= 200) and state.deepVeinEnabled ~= true then
               state.mineDepth = 200
               if state.currentStrata == 0 then
                    send("rendCommandManager",
                         "odinRendererStubMessage",
                         "ui\\orderIcons.xml",
                         "mineshaft_image",
                         "Mine Deepening Halted!", -- header text
                         "Your mineshaft has reached 200 feet and requires a Ventilation Unit to go any deeper!", -- text description
                         "Right-click to dismiss.", -- action string
                         "miningAlert", -- alert type (for stacking)
                         "ui//eventart//mineShoring.png", -- imagename for bg
                         "low", -- importance: low / high / critical
                         state.rOH, -- object ID
                         30 * 1000, -- duration in ms
                         30, -- "snooze" time if triggered multiple times in rapid succession
                         nil) -- gameobjecthandle of director, null if none
               end
			
          elseif (state.mineDepth + digAmount >= 50) and state.coreVeinEnabled ~= true then
               state.mineDepth = 50
               if state.currentStrata == 0 then
                    send("rendCommandManager",
                         "odinRendererStubMessage",
                         "ui\\orderIcons.xml",
                         "mineshaft_image",
                         "Mine Deepening Halted!", -- header text
                         "Your mineshaft has reached 50 feet and requires a Dewatering Pump to go any deeper!", -- text description
                         "Right-click to dismiss.", -- action string
                         "miningAlert", -- alert type (for stacking)
                         "ui//eventart//mineShoring.png", -- imagename for bg
                         "low", -- importance: low / high / critical
                         state.rOH, -- object ID
                         30 * 1000, -- duration in ms
                         30, -- "snooze" time if triggered multiple times in rapid succession
                         nil) -- gameobjecthandle of director, null if none
               end
          else
               state.mineDepth = state.mineDepth + digAmount
          end
		
          --Note: You can add a tech flag for digging more here if you want.
		send("rendUIManager", "SetOfficeString",	SELF, "mineDepth", tostring(state.mineDepth) .. " yards")
          
          if state.currentStrata == 0 then
               checkForButtonEnablement()
               return "productResponse", "none"
          end
          
          --This block adds mining products once you've hit certain depths (and random chance)
          
          if #state.specialProducts == 0 and state.mineDepth > 200 and rand(1,20) == 1 then
			
               -- Reveal thy ore!
               local orePick = ""
               local chance = rand(1,60)
			
               if chance == 60 then
				--You get an SSR product!
                    orePick = state.SSRQualityProductTable[ rand(1,#state.SSRQualityProductTable) ]
               elseif chance > 29 then
				--You get an HQ product!
                    orePick = state.highQualityProductTable[ rand(1,#state.highQualityProductTable) ] 
               else
				--You get a LQ product.
                    orePick = state.lowQualityProductTable[ rand(1,#state.lowQualityProductTable) ]
               end
			
               table.insert(state.specialProducts,orePick)
               
               local headerText = "New Vein Discovered!"
               local bodyText = "Our miners have discovered a new vein of " .. EntityDB[orePick].display_name .. "! \z
                    We will have access to this resource in the Deep Vein stratum."
                    
               if orePick == "logs" then
                    headerText = "Underground Forest Found!"
                    bodyText = "Astonishingly, our miners have discovered a rich forest deep underground - massive mushrooms as far as the eye can see. \z
                    In addition to its previous products, we will now be able to acquire wood from the Deep Vein stratum!"
               elseif orePick == "steel_ingots" then
                    headerText = "Strange Ruins Found!"
                    bodyText = "To the surprise of everyone, our miners have discovered a vast ruin within the Deep Vein Stratum! The ruins appear entirely empty, an labrynthine series of buildings with no explainable purpose. \z
                    While the lack of artefacts disappoints our scientists, the buildings are full of steel materials we could easily repurpose into ingots. An excellent find!"
               elseif orePick == "obeliskian_block" then
                    headerText = "Odd Substance Found!"
                    bodyText = "During an underground expedition, our miners found something most unusual - a strange, pulsating block of stone the like we have never seen before, so large that we cannot find its edge. \z
                    Our scientists are entirely unsure what use such a strange material might have, but a Proper Clockworkian never lets anything like use get in the way of unmitigated harvest! \z
                    Tenebrous Blocks are now available from the Deep Vein stratum!"
               end
               
               send("rendCommandManager",
				"odinRendererTickerMessage",
				"Our miners found a new vein of " .. EntityDB[orePick].display_name .. " !",
				"mineshaft_image",
				"ui\\orderIcons.xml")
               
               send("rendCommandManager",
                    "odinRendererStubMessage",
                    "ui\\orderIcons.xml",
                    "mineshaft_image",
                    headerText, -- header text
                    bodyText, -- text description
                    "Right-click to dismiss.", -- action string
                    "miningAlert", -- alert type (for stacking)
                    "ui//eventart//mineShoring.png", -- imagename for bg
                    "low", -- importance: low / high / critical
                    state.rOH, -- object ID
                    30 * 1000, -- duration in ms
                    30, -- "snooze" time if triggered multiple times in rapid succession
                    nil) -- gameobjecthandle of director, null if none
          end
          
          if #state.specialProducts == 1 and state.mineDepth > 500 and rand(1,100) == 1 then
			
               -- Low chance for a third ore type.
			function giveMeRandomPick()
				local orePick = ""
				local chance = rand(1,60)  --this is to prevent dupe types.
				
                    if chance == 60 then
					--You get an SSR product!
                         orePick = state.SSRQualityProductTable[rand(1,#state.SSRQualityProductTable)]
                    elseif chance > 19 then
					--You get an HQ product!
                         orePick = state.highQualityProductTable[rand(1,#state.highQualityProductTable)] 
                    else
					--You get a LQ product.
                         orePick = state.lowQualityProductTable[rand(1,#state.lowQualityProductTable)]
                    end
				
				return orePick
			end
				
			local orePick = false
			local valid = false
			local count = 0
			while valid == false do
				orePick = giveMeRandomPick()
				count = count + 1
				if count > 20 then
					valid = true
				end
				
				--check to make sure we don't already have this one
				if state.specialProducts[1] ~= orePick then
                         table.insert(state.specialProducts,orePick)
					valid = true
                    end
               end
          
			if orePick ~= false then
				local headerText = "New Ore Discovered!"
				local bodyText = "Our miners have discovered a new vein of " .. EntityDB[orePick].display_name .. "! \z
					We will have access to this resource in the Deep Vein stratum. We are deep enough now that it seems unlikely that we will find further new veins in this mine."
					
				if orePick == "logs" then
					headerText = "Underground Forest Found!"
					bodyText = "Astonishingly, our miners have discovered a rich forest deep underground - massive mushrooms as far as the eye can see. \z
					In addition to its previous products, we will now be able to acquire wood from the Deep Vein stratum!"
					
				elseif orePick == "steel_ingots" then
					headerText = "Strange Ruins Found!"
					bodyText = "To the surprise of everyone, our miners have discovered a vast ruin within the Deep Vein Stratum! The ruins appear entirely empty, an labrynthine series of buildings with no explainable purpose. \z
					While the lack of artefacts disappoints our scientists, the buildings are full of steel materials we could easily repurpose into ingots. An excellent find!"
					
				elseif orePick == "obeliskian_block" then
					headerText = "Odd Substance Found!"
					bodyText = "During an underground expedition, our miners found something most strange - a strange, pulsating block of stone the like we have never seen before, so large that we cannot find its edge. \z
					Our scientists are entirely unsure what use such a strange material might have, but a Proper Clockworkian never lets anything like use get in the way of unmitigated harvest! \z
					Tenebrous Blocks are now available from the Deep Vein stratum!"
					
				end
				
				send("rendCommandManager",
					"odinRendererStubMessage",
					"ui\\orderIcons.xml",
					"mineshaft_image",
					headerText, -- header text
					bodyText, -- text description
					"Right-click to dismiss.", -- action string
					"miningAlert", -- alert type (for stacking)
					"ui//eventart//mineShoring.png", -- imagename for bg
					"low", -- importance: low / high / critical
					state.rOH, -- object ID
					30 * 1000, -- duration in ms
					30, -- "snooze" time if triggered multiple times in rapid succession
					nil) -- gameobjecthandle of director, null if none
				
			end
          end
		
		UpdateMineProducts()
		
          -- End product block
          
          checkForButtonEnablement() --see if there's any new buttons to enable while we're at it
          
          local productTable = {}
          for k,v in pairs(state.defaultProducts["strata" .. tostring(state.currentStrata)]) do
               table.insert(productTable,v)
          end
		
          if state.currentStrata > 1 then -- add ore products
               for k,v in pairs(state.localProducts) do
                    if v == true then
                         table.insert(productTable,k)
                    end
               end
          end
		
          if state.currentStrata > 3 then -- add 2nd ore vein
               for k,v in pairs(state.specialProducts) do
                    table.insert(productTable,v)
               end
          end
		
          if state.currentStrata > 4 and rand(1,15) == 1 then -- add artifacts, but only rarely.
               for k,v in pairs(state.rareItems) do
                    table.insert(productTable,v)
               end
          end
          
          if query("gameSession", "getSessionBool", "digging1_unlocked")[1] == true then
			--Do you have mining tech that improves chances of specific things?
			if rand(1,4) + state.currentStrata >= 6 then
				--Basically, finds all ores in the list and picks one with a chance of 1/4 at type 2 and 1/1 at type 5.
				local addTable = {}
                    for k,v in pairs(productTable) do
					if v == "hematite" or
						v == "malachite" or
						v == "sphalerite" or
						v == "native_gold" then
						
						table.insert( addTable, v)
					end
				end
                    for k,v in pairs(addTable) do --this is to avoid embarrassing infinite loops.
                         table.insert( productTable, v)
                    end
			end
		end
		
		local product = productTable[rand(1,#productTable)]
		
		-- update artifact stats
		local product_tags = EntityDB[ product ].tags
		for k,v in pairs(product_tags) do
			if v == "artifact" then
				send("gameSession","incSessionInt","eldritchArtifactsFound", 1)
				send("gameSession","setSessionString","endGameString10",
					tostring( query("gameSession","getSessionInt","eldritchArtifactsFound")[1]) )
				
				break
			end
		end
			
		if state.mineDepth > 500 then
			if not query("gameSession","getSessionBool","delveDeep")[1] then
				send("gameSession","setSessionBool","delveDeep",true)
				send("gameSession", "setSteamAchievement", "delveDeep")
			end
		end
			
          return "productResponse", product
	>>
     
     receive AddLocalProduct(string product)
	<<
          state.localProducts[product] = true
	>>
	
	receive RemoveLocalProduct(string product)
	<<
          state.localProducts.product = nil
	>>
     
     receive RefreshLocalProducts()
     <<
          local modules = query("gameWorkshopManager", "getBuildingModulesGameSide", SELF)[1]
		
          for k,v in pairs(modules) do
               if v.tags.mineshaft then 
                    send(v,"getProducts")
               end
          end
		
          UpdateMineProducts()
     >>
     
     receive PurgeLocalProducts()
     <<
          for k,v in pairs(state.localProducts) do
               state.localProducts[k] = nil
          end
		
          send (SELF, "RefreshLocalProducts")
     >>

	respond GetMineDepth()
	<<
		return "MineDepthMessage", state.mineDepth
	>>
     
     respond GetMineStrata()
	<<
		return "StrataMessage", state.currentStrata
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
		mine_reset_supply_text()
	>>
	
	receive consumeSupplies( int tier, int count )
	<<
		state.supplies[tier] = state.supplies[tier] - count 
		send("rendUIManager", "SetOfficeInt", SELF, "workPoints" .. tostring(tier) , state.supplies[tier])
		mine_reset_supply_text()
	>>
	
	receive refreshTechModifierDisplay()
	<<
		local techmod = "N/A"
		if state.parent.techName then
			techmod = 100 - query("gameSession","getSessionInt", state.parent.techName )[1]
			techmod = "+" .. tostring( techmod ) .. "%"
		end
		
		send("rendUIManager",
			"SetOfficeString",
			SELF,
			"workshopTechModifier",
			techmod )
	>>
	
	receive setBuildingOwner( gameObjectHandle newOwner )
	<<
		if newOwner and state.supplies[1] == 0 then
			
			local eventQ = query("gameSimEventManager",
								"startEvent",
								"supplies_warning_mine",
								{},
								{} )[1]
						
			send(eventQ,"registerBuilding",SELF)
		end
	>>
>>