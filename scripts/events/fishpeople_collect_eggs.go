event "fishpeople_collect_eggs"
<<
	state 
	<<
	>>

	receive Create( stringstringMapHandle name )
	<< 
		printl("events", "Fishpeople egg collection event started")
	>>
	
	FSM 
	<<
		["start"] = function(state,tags)
			local targetEggs = tags["fishpeople_eggs"].target
			local targetPos = query(targetEggs,"gridGetPosition")[1]
		
			settimer("Fishpeople Event Timer", 0)				
               
			local s = "an unidentified group"
			local icon = "mysterious_figures"
			local fishSeen = query("gameSession", "getSessionBool", "fishpeopleFirstContact")[1]
			if fishSeen then
				s = "a group of Fishpeople"
				icon = "fishperson"
			end
			
			send("rendCommandManager", "odinRendererTickerMessage",
				"The Imperial Airship Corps reports that " .. s .. " has been spotted.",
                    icon,
                    "ui\\thoughtIcons.xml")
			
			send("rendCommandManager",
				"odinRendererStubMessage",
					"ui\\thoughtIcons.xml", -- iconskin
					"fishperson", -- icon
					"Fishpeople Spotted", -- header text
					"The Imperial Airship Corps reports that " .. s .. " has been spotted.", -- text description
					"Right-click to dismiss.", -- action string
					"fishpeople_eggs", -- alert type (for stacking)
					"ui//eventart//fishPeople.png", -- imagename for bg
					"low", -- importance: low / high / critical
					nil, -- object ID
					45 * 1000, -- duration in ms
					0, -- "snooze" time if triggered multiple times in rapid succession
					nullHandle)
			
			local spawnTable = { legacyString = "Fishy Patrol Group" }
               local fish_group = query("scriptManager",
								"scriptCreateGameObjectRequest",
								"fishpeople_group",
								spawnTable)[1]
			
               send(fish_group,"GameObjectPlace", -1, -1 )
               send(fish_group,"pushMission","collect_eggs",targetPos, 1 )
			
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
