event "fishpeople_nesting"
<<
	state 
	<<
	>>

	receive Create( stringstringMapHandle name )
	<< 
		printl("events","Fishpeople nesting event started")
	>>

	FSM 
	<<
		["start"] = function(state,tags)
			settimer("Fishpeople Event Timer", 0)				
               
			-- for some mystery
			if rand(1,2) == 1 then
				local s = "an unidentified group"
				local icon = "mysterious_figures"
				local t = "Unidentified"
				
				local fishSeen = query("gameSession", "getSessionBool", "fishpeopleFirstContact")[1]
				if fishSeen then
					-- even more mystery
					if rand(1,2) == 1 then
						s = "a group of Fishpeople"
						icon = "fishperson"
						t = "Fishpeople"
					end
				end
				
				send("rendCommandManager", "odinRendererTickerMessage",
					"The Imperial Airship Corps reports that " .. s .. " has been spotted.",
					icon,
					"ui\\thoughtIcons.xml")
                    send("rendCommandManager",
                         "odinRendererPlaySoundMessage",
                         "alertNeutral")
						 
				send("rendCommandManager",
					"odinRendererStubMessage",
					"ui\\thoughtIcons.xml", -- iconskin
					icon, -- icon
					"" .. t .. " Gathering", -- header text
					"The Imperial Airship Corps reports that " .. s .. " has been spotted.", -- text description
					"Right-click to dismiss.", -- action string
					"fishperson_noninteractive", -- alert type (for stacking)
					"", -- imagename for bg
					"low", -- importance: low / high / critical
					nil, -- object ID
					45 * 1000, -- duration in ms
					0, -- "snooze" time if triggered multiple times in rapid succession
					nullHandle)
			end

			local spawnTable = { legacyString = "Fishy Patrol Group" }
               local fish_group = query( "scriptManager",
								"scriptCreateGameObjectRequest",
								"fishpeople_group",
								spawnTable)[1]
			
               send(fish_group,"GameObjectPlace", -1, -1 )
               send(fish_group,"pushMission","make_nest",nil, 5)
               
			
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
