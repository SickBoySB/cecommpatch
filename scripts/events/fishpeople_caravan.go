event "fishpeople_caravan"
<<
	state 
	<<
	>>

	receive Create( stringstringMapHandle name )
	<< 
		printl("events","Fishpeople convoy event started")
	>>

	FSM 
	<<
		["start"] = function(state,tags)
			settimer("Fishpeople Event Timer", 0)				
               
			if query("gameSession", "getSessionBool", "biomeDesert")[1] then
				if rand(1,2) == 1 then
					printl("events", "hit 50% chance to abort fishpeople event due to desert")
					return "final", true
				end
			end
			
			-- Let's add some mystery.
			-- EXPERIMENTAL: adding more randomness and mystery
			--if rand(1,2) == 1 then
				local s = "an unidentified convoy"
				local icon = "mysterious_figures"
				local fishSeen = query("gameSession", "getSessionBool", "fishpeopleFirstContact")[1]
				if fishSeen then
					-- even more mystery
					if rand(1,2) == 1 then
						s = "a convoy of Fishpeople"
						icon = "fishperson"
					end
				end
				
				send("rendCommandManager",
				"odinRendererStubMessage",
				"ui\\thoughtIcons.xml", -- iconskin
				icon, -- icon
				"Unusual Caravan", -- header text
				"Recent reports from passing airships indicate that " .. s .. " carrying goods has been spotted walking toward the sea in the vicinity of your colony.", -- text description
				"Right-click to dismiss.", -- action string
				"fishperson_noninteractive", -- alert type (for stacking)
				"", -- imagename for bg
				"low", -- importance: low / high / critical
				nil, -- object ID
				30 * 1000, -- duration in ms
				0, -- "snooze" time if triggered multiple times in rapid succession
				nil)
                    
                    send("rendCommandManager",
                         "odinRendererPlaySoundMessage",
                         "alertNeutral")
				
			--end
			
			local spawnTable = { legacyString = "Fishy Patrol Group" }
			if rand(1,2) == 1 then
				spawnTable.bringGifts = "true" -- intentional string
			else
				spawnTable.bringFood = "true" -- intentional string
			end
			
               local fish_group = query("scriptManager",
								"scriptCreateGameObjectRequest",
								"fishpeople_group",
								spawnTable)[1]
			
               send(fish_group,"GameObjectPlace", -1, -1)
               send(fish_group,"pushMission","patrol",nil, 1 )

			return "final"
		end,

		["final"] = function(state,tags)  	
			return
		end,

		["abort"] = function(state,tags) 
			return
		end
	>>
>>
