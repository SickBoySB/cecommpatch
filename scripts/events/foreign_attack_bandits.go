event "foreign_attack_bandits"
<<
	state 
	<<
		bool alertTriggered
          int counter
		string dialogBoxResults
		int timeout
		gameGridPosition targetPos
	>>

	receive Create( stringstringMapHandle name )
	<< 
		printl("events","Foreign attack bandits event started")
	>>

	receive boxSelectionResult( string id, string message )
	<<
		state.dialogBoxResults = message
	>>

	receive alertSelectedResult( string message )
	<<
		state.alertTriggered = true
	>>
	
	FSM 
	<<
		["start"] = function(state,tags)
			settimer("Foreign Troops Event Timer", 0)
			
			local dominant = query("gameSession","getSessionString", "dominantFaction")[1]
			
			-- CECOMMPATCH more effing hacks for dominant not using the right name
			local dominant_short_name = {
				["Grossherzoginnentum von Stahlmark"] = "Stahlmark",
				["Novorus Imperiya"] = "Novorus",
				["Republique Mecanique"] = "Republique",
				-- not sure if any of the below is ever used, but just to be safe..
				["The Clockwork Empire"] = "Empire",
				["Fishpeople"] = "Fishpeople",
				["Bandits"] = "Bandits"
			}
				
			if dominant_short_name[dominant] then
				--printl("CECOMMPATCH - dominant faction shortname correction")
				dominant = dominant_short_name[dominant]
			end
			-- /hack
			
			local nations = {"Stahlmark", "Novorus", "Republique", dominant }
			local nation = nations[ rand(1, #nations) ]
	
			local isPatrol = false
			if rand(1,2) == 1 then isPatrol = true end
			local isNeutral = true
			local isAllied = query("gameSession", "getSessionBool", nation .. "Friendly")[1] 
			local isHostile = query("gameSession", "getSessionBool", nation .. "Hostile")[1]
			if isAllied or isHostile then
				isNeutral = false
			end
			
			local nationInfo = EntityDB[ nation .. "Info"]
			
               state.alertTriggered = false
			local s = "A unit of " .. nationInfo.adjective .. " troops is attacking a local bandit camp."
			
			if isAllied then
				s = "A unit of allied " .. nationInfo.adjective .. " is attacking a local bandit camp."
			elseif isHostile then
				s = "A unit of hostile " .. nationInfo.adjective .. " is attacking a local bandit camp."
			else
				s = "A unit of neutral " .. nationInfo.adjective .. " is attacking a local bandit camp."
			end
			
			-- endLoc is the position of the targeted bandit.
			local bandit = tags["character"].target
			local endLoc = query(bandit, "gridGetPosition")[1]
			
			printl("events", "Foreigner squad targeting bandit at: " .. tostring(endLoc.x) .. " / " .. tostring(endLoc.y) )
			
			if endLoc == false or
				startLoc == false then
				
				return "abort", true
			end

               local foreign_group = query( "scriptManager",
									"scriptCreateGameObjectRequest",
									"foreigner_group",
									{ legacyString = nation .." Unit",
									mission = "idle", } )[1]
			
               send(foreign_group,"GameObjectPlace", -1 , -1 )

               send(foreign_group,"pushMission", "patrol", endLoc, 15)

			-- Okay, NOW tell them we did it, after successful placement
			send("rendCommandManager",
				"odinRendererTickerMessage",
				s,
                    nationInfo.iconTroops,
                    "ui\\thoughtIcons.xml")
			
			local isAllied = query("gameSession", "getSessionBool", nation .. "Friendly")[1] 
			if isAllied then
				
				local leader = query(foreign_group,"getLeader")[1]
				local leaderhandle = query(leader,"ROHQueryRequest")[1]
				
				send("rendCommandManager",
					"odinRendererStubMessage",
					"ui\\thoughtIcons.xml", -- iconskin
					nationInfo.iconTroops, -- icon
					"Allies attack Bandits", -- header text
					s, -- text description
					"Left-click to zoom. Right-click to dismiss.", -- action string
					"foreignAttackBandits", -- alert type (for stacking)
					"", -- imagename for bg
					"low", -- importance: low / high / critical
					leaderhandle, -- object ID
					45 * 1000, -- duration in ms
					0, -- "snooze" time if triggered multiple times in rapid succession
					nullHandle)
				
			end
			
			return "final", true
		end,

		["final"] = function(state,tags)  
			return
		end,

		["abort"] = function(state,tags)  
			return
		end
	>>
>>
