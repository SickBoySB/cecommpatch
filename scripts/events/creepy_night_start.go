event "creepy_night_start"
<<
	state 
	<<
		bool alertTriggered
          int counter
		string dialogBoxResults
		bool lock
	>>

	receive Create( stringstringMapHandle init )
	<< 
		printl("events","creepy_night_start started")
		state.director = query("gameSession","getSessiongOH","event_director" .. init.director_name)[1]
	>>
	
	FSM 
	<<
		["start"] = function(state,tags)
			send("rendCommandManager",
				"odinRendererFYIMessage",
				"ui\\thoughtIcons.xml", -- iconskin
				"mad", -- icon
				"Unnatural Sounds from the Night", -- header text
				"As the sun sets, strange howls begin beyond the outskirts of town, invoking imagined figures of unnatural animals being held at bay by the glimmer of the settlement lights.",
				"(Left-Click to read more, Right-click to clear)", -- tooltip string 
				"creepynight", -- alert type (for stacking)
				"ui//eventart//strange_lights_forest.png", -- imagename for bg
				"high", -- importance: low / high / critical
				nil, -- state.renderHandle, -- object ID
				180 * 1000, -- duration in ms
				0, -- "snooze" time if triggered multiple times in rapid succession
				state.director)
				--memories
				
				local citizens = query("gameSpatialDictionary", "allCharactersWithTagRequest", "citizen")
				
				if citizens and citizens[1] then
					for k,v in pairs(citizens[1]) do
						local tags = query(v,"getTags")[1]
						if not tags.dead then
							--for k,v in pairs(citizens[1]) do
								send(v,"makeMemory","Unnerved by Howling in the Night",nil, nil, nil, nil)
							--end
						end
					end
				end
				
                    send("rendCommandManager",
					"odinRendererPlaySoundMessage",
					"alertNeutral")
				
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
