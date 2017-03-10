event "fishpeople_patrol"
<<
	state 
	<<
          int counter
          table spawnLocation
		string dialogBoxResults
		int timeout
	>>

	receive Create( stringstringMapHandle name )
	<< 
		printl("events","Fishpeople beach patrol event started")
	>>

	receive boxSelectionResult( string id, string message )
	<<
		state.dialogBoxResults = message
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
			
			local s = "an unidentified group"
			local icon = "mysterious_figures"
			local fishSeen = query("gameSession", "getSessionBool", "fishpeopleFirstContact")[1]
			if fishSeen then
				s = "a group of Fishpeople"
				icon = "fishperson"
			end
			
			send("rendCommandManager",
				"odinRendererStubMessage",
				"ui\\thoughtIcons.xml", -- iconskin
				icon, -- icon
				"Fishpeople Patrol", -- header text
				"Recent reports from passing airships indicate that " .. s .. " has been spotted patrolling near your colony. \z
					Best keep an eye out, and your miltiary well-staffed.", -- text description
				"Right-click to dismiss.", -- action string
				"fishperson_noninteractive", -- alert type (for stacking)
				"", -- imagename for bg
				"low", -- importance: low / high / critical
				nil, -- object ID
				45 * 1000, -- duration in ms
				0, -- "snooze" time if triggered multiple times in rapid succession
				nullHandle)
               
               send("rendCommandManager",
                         "odinRendererPlaySoundMessage",
                         "alertNeutral")
               
			
               local fish_group = query( "scriptManager",
								"scriptCreateGameObjectRequest",
								"fishpeople_group",
								{ legacyString = "Fishy Patrol Group", } )[1]
			
               send(fish_group, "GameObjectPlace", -1, -1)
               send(fish_group,"pushMission","patrol", nil, 15 )

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
