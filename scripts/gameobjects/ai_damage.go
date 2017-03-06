gameobject "ai_damage"
<<
	local
	<<
		function doSFXandPFX( damagingObject, damageType, damageAmount, onhit_effect )
			
			local sfxName = false
			local pfxName = false
			
			if damageType == "squish" then
				sfxName = "Squish"
			elseif damageType == "punch" or damageType == "blunt" then
				if SELF.tags.armoured and not SELF.tags.animal and not SELF.tags.obeliskian then
					sfxName = "Metal Plate"
					-- sparks?
				elseif SELF.tags.armoured and (SELF.tags.animal or SELF.tags.obeliskian) then
					sfxName = "Stone"
					-- dust?
				else 
					sfxName = "Punch"
					pfxName = "BloodSplashCentered"
				end
			elseif damageType == "slash" then
				if SELF.tags.armoured then
					sfxName = "Bullet Ricochet"
				else 
					sfxName = "Slice Flesh"
					pfxName = "BloodSplashCentered"
				end
			elseif damageType == "pierce" then
				if SELF.tags.armoured then
					sfxName = "Bullet Ricochet"
				else 
					sfxName = "Flesh Pierce"
					pfxName = "BloodSplashCentered"
				end
				
			elseif damageType == "bullet" then
				if SELF.tags.armoured then
					sfxName = "Bullet Ricochet"
				else 
					sfxName = "Bullet Hit"
					pfxName = "BloodSplashCentered"
				end
				
			elseif damageType == "explosion" then
				if not SELF.tags.armoured then
					pfxName = "BloodSplashCentered"
				end
				
			elseif damageType == "mind_blast" then
				-- what was I doing again?
				-- maybe make people reel with psychic terror
				if not SELF.tags.obeliskian and not SELF.tags.spectre then
					--send(SELF,"AICancelJob", "hit by psychic blast")
				end
				
				if SELF.tags.citizen then
					send(SELF,"makeMemory","Struck By Mind Blast",nil,nil,nil,nil)
				end
				pfxName = "QuagSmokePuff"
				
			elseif damageType == "eldritch" then
				pfxName = "BloodSplashCentered"
				
			elseif damageType == "fire" then
				pfxName = "SmallSplosion"
				
				-- to stop inevitable death by fire
				if rand(1,10) == 1 then
					send(SELF,"putOutFire")
				end
			elseif damageType == "voltaic" then
				-- pfx:	"ElecSparker"
				sfxName = "Leyden Hit"
			else	
				-- or, well, just do some damage.
				pfxName = "BloodSplashCentered"
			end
			
			if sfxName then
				send("rendInteractiveObjectClassHandler",
					"odinRendererPlaySFXOnInteractive",
					state.renderHandle,
					sfxName)
			end
			
			if pfxName then
				send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					pfxName,
					state.AI.position.x,
					state.AI.position.y)
			end	
		end
		
		function doOnHitEffects( damagingObject, damageType, damageAmount, onhit_effect )
			
			if onhit_effect == "set_on_fire" then
				send(SELF, "IgniteMessage")
				
			elseif onhit_effect == "fire_blast" then
				
				send("rendInteractiveObjectClassHandler",
					"odinRendererPlaySFXOnInteractive",
					state.renderHandle,
					"Explosion Generic")
				
				send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					"MidSplosion",
					state.AI.position.x,
					state.AI.position.y)
	
				local results = query("gameSpatialDictionary", "allObjectsInRadiusRequest", state.AI.position, 1, false)
				if results ~= nil then
					send(results[1], "damageMessage", damagingObject, "fire", 3, "set_on_fire"  )
				end
				
			elseif onhit_effect == "Quaggaroth Trample" then
				local handle = query("scriptManager",
								"scriptCreateGameObjectRequest",
								"explosion",
								{ legacyString= "Quaggaroth Trample",
								 x = tostring(state.AI.position.x),
								 y = tostring(state.AI.position.y), } )[1]
									
			elseif onhit_effect == "Quaggaroth Beam Hit" then
				local handle = query("scriptManager",
								"scriptCreateGameObjectRequest",
								"explosion",
								{ legacyString= "Quaggaroth Beam Hit",
								 x = tostring(state.AI.position.x),
								 y = tostring(state.AI.position.y), } )[1]
				
			elseif onhit_effect == "voltaic_hit_small" then
				if rand(1,2) == 1 then send(SELF, "IgniteMessage") end
			elseif onhit_effect == "voltaic_hit_medium" then
				if rand(1,2) == 1 then send(SELF, "IgniteMessage") end
				local transforms = {
						{0,1},
				  {-1,0}, {0,0}, {1,0},
					    {0,-1}, }
				
				for k,v in pairs(transforms) do
					send("rendCommandManager",
						"odinRendererCreateParticleSystemMessage",
						"ElecSparker",
						state.AI.position.x + v[1],
						state.AI.position.y + v[2])
				end
					
			elseif onhit_effect == "voltaic_hit_large" then
				if rand(1,2) == 1 then send(SELF, "IgniteMessage") end
				local transforms = {
							{0,2},
					  {-1,1}, {0,1}, {1,1},
				{-2,0},{-1,0}, {0,0}, {1,0}, {2,0},
					  {-1,-1},{0,-1},{1,-1},
							{0,-2},  }
				
					for k,v in pairs(transforms) do
						send("rendCommandManager",
							"odinRendererCreateParticleSystemMessage",
							"ElecSparker",
							state.AI.position.x + v[1],
							state.AI.position.y + v[2])
					end
					
			elseif onhit_effect == "Blunderbuss Area Attack" then
				handle = query( "scriptManager",
							 "scriptCreateGameObjectRequest",
							 "explosion",
							 { legacyString = "Blunderbuss Area Attack" } )[1]
				
				send( handle,
					"GameObjectPlace",
					state.AI.position.x,
					state.AI.position.y  )
				
				
			elseif onhit_effect == "Small Eldritch Explosion" then
				handle = query( "scriptManager",
							 "scriptCreateGameObjectRequest",
							 "explosion",
							 { legacyString = "Small Eldritch Explosion" } )[1]
				
				send( handle,
					"GameObjectPlace",
					state.AI.position.x,
					state.AI.position.y  )
				
			elseif onhit_effect == "Psychic Blast" then
				local results = query("gameSpatialDictionary","allObjectsInRadiusRequest", state.AI.position, 5, false)
				if results ~= nil then
					send(results[1], "damageMessage", damagingObject, "mind_blast", 5, ""  )
				end
				
			elseif onhit_effect == "Small Psychic Blast" then
				local results = query("gameSpatialDictionary","allObjectsInRadiusRequest", state.AI.position, 2, false)
				if results ~= nil then
					send(results[1], "damageMessage", damagingObject, "mind_blast", 2, ""  )
				end
				
			elseif onhit_effect == "Minum Burst Hit" then
				
				local results = query( "scriptManager",
							 "scriptCreateGameObjectRequest",
								"explosion",
								{ legacyString= onhit_effect } )
				if results then
					if results[1] then
						send( results[1], "GameObjectPlace", state.AI.position.x,state.AI.position.y )
					end
				end
				
			elseif onhit_effect == "Spawn Primed Grenade Cluster" then
				
				for i=1,3 do 
					if rand(1,2) == 1 then
						local handle = query("scriptManager",
										"scriptCreateGameObjectRequest",
										"clearable",
										{legacyString= "Primed Grenade",
										 timer = tostring(rand(1,21)) })[1]
							
						send(handle,
							"GameObjectPlace",
							state.AI.position.x + rand(-1,1),
							state.AI.position.y + rand(-1,1))
					else
						-- explode immediately.
						local handle = query("scriptManager",
										"scriptCreateGameObjectRequest",
										"explosion",
										{legacyString= "Medium Explosion" })[1]
						
						send(handle,
							"GameObjectPlace",
							state.AI.position.x + rand(-1,1),
							state.AI.position.y + rand(-1,1))
						
					end
				end
				
			elseif onhit_effect == "Spawn Primed Grenade" then
				
				if rand(1,2) == 1 then
					local handle = query("scriptManager",
									"scriptCreateGameObjectRequest",
									"clearable",
									{legacyString= "Primed Grenade",
									 timer = tostring(rand(1,21)) })[1]
						
					send(handle,
						"GameObjectPlace",
						state.AI.position.x,
						state.AI.position.y)
				else
					-- explode immediately.
					local handle = query("scriptManager",
									"scriptCreateGameObjectRequest",
									"explosion",
									{legacyString= "Medium Explosion" })[1]
					
					send(handle,
						"GameObjectPlace",
						state.AI.position.x,
						state.AI.position.y)
					
				end
				
			elseif onhit_effect == "cannon_hit" then
				
				local handle = query("scriptManager",
									"scriptCreateGameObjectRequest",
									"explosion",
									{legacyString= "Medium Explosion" })[1]
					
					send(handle,
						"GameObjectPlace",
						state.AI.position.x,
						state.AI.position.y)
					
				send("rendInteractiveObjectClassHandler",
					"odinRendererPlaySFXOnInteractive",
					state.renderHandle,
					"Large Bomb")
				
				send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					"DustExplosion",
					state.AI.position.x,
					state.AI.position.y)
				
			elseif onhit_effect == "grenade_explosion" then
	
				send("rendInteractiveObjectClassHandler",
					"odinRendererPlaySFXOnInteractive",
					state.renderHandle,
					"Explosion Generic")
				
				send("rendCommandManager",
					"odinRendererCreateParticleSystemMessage",
					"MidSplosion",
					state.AI.position.x,
					state.AI.position.y)
	
				local results = query("gameSpatialDictionary", "allObjectsInRadiusRequest", state.AI.position, 3, false)
				if results ~= nil then
					send(results[1], "damageMessage", damagingObject, "blunt", 2, ""  )
					send(results[1], "damageMessage", damagingObject, "shrapnel", 3, ""  )
				end
	
				-- and then place a fiery crater because it's fun
				local results = query( "scriptManager", "scriptCreateGameObjectRequest", "clearable", { legacyString="Medium Crater" } )
				if results then
					send( results[1], "GameObjectPlace", state.AI.position.x,state.AI.position.y )
				end
				
			elseif onhit_effect == "Urchin Grenade Explosion" then
				local results = query( "scriptManager",
							 "scriptCreateGameObjectRequest",
								"explosion",
								{ legacyString= onhit_effect } )
				if results then
					if results[1] then
						send( results[1], "GameObjectPlace", state.AI.position.x,state.AI.position.y )
					end
				end
			end
		end
	>>
	
	state
	<<
		table myAfflictions
	>>

	receive Create( stringstringMapHandle init )
	<<
		state.myAfflictions = {}
		SELF.tags.ignite_128 = true
	>>

	receive damageMessage( gameObjectHandle damagingObject, string damageType, int damageAmount, string onhit_effect )
	<<
		if SELF.deleted then
			return
		end
		
		if not state.AI then
			return
		end
		
		if SELF.tags["dead"] then
			if state and state.AI and not SELF.deleted then
				doSFXandPFX( damagingObject, damageType, damageAmount, onhit_effect )
				doOnHitEffects( damagingObject, damageType, damageAmount, onhit_effect )
				
				if SELF.tags.horror_corpse then
					state.AI.ints["health"] = state.AI.ints["health"] - damageAmount
					if state.AI.ints.health < -30 then
						
						send("rendCommandManager",
							"odinRendererCreateParticleSystemMessage",
							"QuagSmokePuff",
							state.AI.position.x,
							state.AI.position.y)
						
						send(SELF,"HarvestMessageNoProduct",SELF,nil)
					end
				end
				return
			else
				-- you don't really exist; weird!
				return
			end
		end
		
		if damageType == "blunderbuss" and SELF.tags.citizen then
			-- friendlies will ignore blunderbuss shrapnel damage
			-- they DO take "explosion" damage however, but its lower and more localized.
			return
		end
		
		-- DO damage bonus from military training & brutish trait
		
		local damagingObjectTags = {}
		local damageBonus = 0
		
		if damagingObject then 
			damagingObjectTags = query(damagingObject, "getTags")
			if damagingObjectTags and damagingObjectTags[1]["citizen"] then
				local results = query(damagingObject, "getAIAttributes")
				if results then
					local otherAI = results[1]
					if otherAI then
						if query(damagingObject,"hasTrait","Brutish")[1] then
							-- you monster
							damageBonus = damageBonus + 1
						end

						if damagingObjectTags[1]["lower_class"] then
							-- if redcoat, get bonus.
							if damagingObjectTags.military and
								not damagingObjectTags.militia then
								damageBonus = damageBonus + 1
							end
							
							-- if overseer has skill do bonus damage
							local overseer = query("gameBlackboard",
											   "gameObjectGetOverseerMessage",
											   otherAI.currentWorkParty)
							
							if overseer then
								if overseer[1] then
									local overseerState = query(overseer[1], "getAIAttributes")[1]
									if overseerState then
										if overseerState.skills.militarySkill > 1 then
											damageBonus = damageBonus + overseerState.skills.militarySkill -1  
										end
									end
								end
							end
						elseif damagingObjectTags[1]["middle_class"] then
							-- or I am the overseer!
							if otherAI.skills.militarySkill > 1 then
								damageBonus = damageBonus + otherAI.skills.militarySkill -1  
							end
						end
						
						-- do tech modification.
						damageBonus = damageBonus + query("gameSession","getSessionInt","militaryDamageTechBonus")[1]
					end
				end
			end
		end
		
		damageAmount = damageAmount + damageBonus
		
		-- MC chars take slightly less damage due to victorian-era classism
		if SELF.tags.citizen and SELF.tags.middle_class then
			if damageAmount > 1 then
				damageAmount = damageAmount - 1
			end
		end
		
		if SELF.tags.citizen and SELF.tags.military then
			damageAmount = damageAmount - query("gameSession","getSessionInt","militaryDefenseTechBonus")[1]
		end
		
		if damageType == "fire" and SELF.tags.fire_immune then
			damageAmount = 0
		end

		-- TODO: logic for death by damagetype? Maybe do that special-case per ai entity.\

		local armourReductionFactor = 2
		if SELF.tags.spectre then
			-- nope!
			
		elseif SELF.tags.vulnerable then
			state.AI.ints["health"] = state.AI.ints["health"] - (damageAmount * 1) --Can make it 2x or worse later if we want. For now, overrides armor.
		elseif SELF.tags.selenian and
			damageType ~= "fire" and
			damageType ~= "eldritch" then
			
			state.AI.ints["health"] = state.AI.ints["health"] - 1
					
		elseif SELF.tags["armoured"] and
			damageType ~= "fire" and
			damageType ~= "eldritch" then
	
			if SELF.tags.quaggaroth then
				if SELF.tags.vulnerable then
					state.AI.ints["health"] = state.AI.ints["health"] - div(damageAmount, armourReductionFactor)
				else
					state.AI.ints["health"] = state.AI.ints["health"] - 1
				end
			else
				state.AI.ints["health"] = state.AI.ints["health"] - div(damageAmount, armourReductionFactor)
			end
		else
			state.AI.ints["health"] = state.AI.ints["health"] - damageAmount
		end

		-- do sfx and pfx
		doSFXandPFX( damagingObject, damageType, damageAmount, onhit_effect )
		
		-- and now, the additional Fun effects.
		doOnHitEffects( damagingObject, damageType, damageAmount, onhit_effect )
		
		send(SELF, "InCombat")

		if not SELF.tags.spectre and state.AI.ints["health"] <= 0 then

			--printl(state.AI.name .. " taking an affliction");
			-- DAMAGE CAN ONLY CAUSE ONE AFFLICTION AT A TIME
			-- NUM OF AFFLICTIONS NOT NEEDED IN AN INT - WE HAVE A TABLE SIZE
			-- AFFLICTIONS SHOULD BE ADDED TO A TABLE OWNED BY THE AI_DAMAGE OBJECT

			if state.AI.curJobInstance then
				FSM.abort(state, query(SELF,"getName")[1] .. " FSM aborted due to taking affliction!" )
			end

			if SELF.tags["fishperson"] then
				state.AI.ints["numAfflictions"] = state.AI.ints["numAfflictions"] + 1
			elseif SELF.tags["human"] then
				state.AI.ints["numAfflictions"] = state.AI.ints["numAfflictions"] + 1
			end

			send( SELF, "createAffliction", damagingObject, damageType )
		end

	>>

	receive createAffliction( gameObjectHandle damagingObject, string damageType ) -- int overkill, string damageType )
	<<
		--SPAWN GIBS					No.  This should be a message sent to self that's specific to what it is.  People, fish, animals should have different ones.
		--ADD AFFLICTION				Yep, we should do this.  No, we should not have them as integers.
		--CHECK IF WE SHOULD BE DEAD	Yep, we should do this.
		--IF SO, DIE. ...				No.  Need to use a unified "death" message.  Since everything that can die should have ai_damage, do generic stuff here, message for specifics
		--ADD MEMORY  ...				Yes, but don't need to test for human tag.  Things without the message receiver don't care
		--RETURN TRUE  ...				No.  This is unnecessary.  Should be a receive, not a respond.

		-- NO NEED FOR THIS TO BE A MESSAGE.  ADD AFFLICTION, SPAWN GIBS, CHECK FOR DEATH, DIE, AND ADD MEMORY IN THE DAMAGE MESSAGE
		
		if SELF.tags.marked_for_beating then
			-- After getting an affliction, the beatings can end.
			SELF.tags.marked_for_beating = nil
		end
		
		send( SELF, "spawnGibs" )

		-- Count possible afflictions
		local possibles = 0

		local resultingAffliction = nil

		local afflictions = EntityDB["afflictionsDB"].afflictions

		for affName, affliction in pairs(afflictions) do		  
			local matchType = false
			local beingMatch = false
			
			for j, damage in ipairs(affliction.validDamageTypes) do
				if (damage == damageType) or (damage == "all") then
					matchType = true
					break
				end
			end

			if matchType then
				for j, being in ipairs(affliction.validBeings) do
					if SELF.tags[being] or (being == "all") then
						beingMatch = true
						break
					end
				end
				if beingMatch then
					possibles = possibles + 1
				end
			end
		end

		local result = rand(1, possibles)
		local curResult = 1

		for affName, affliction in pairs(afflictions) do
			local matchType = false
			local beingMatch = false
			
			for j, damage in ipairs(affliction.validDamageTypes) do
				if (damage == damageType) or (damage == "all") then
					matchType = true
					break
				end
			end

			if matchType then
				for j, being in ipairs(affliction.validBeings) do
					if SELF.tags[being] or (being == "all") then
						beingMatch = true
						break
					end
				end

				if beingMatch then
					if curResult == result then
						-- We've found our affliction!
						resultingAffliction = affliction
						break
					end
					curResult = curResult + 1
				end
			end
		end		  

		if resultingAffliction then
			local affIcon = "affliction"

			if resultingAffliction.icon then
				affIcon = resultingAffliction.icon
			end

			if state.myAfflictions and state.myAfflictions[1] then
				state.myAfflictions[ #state.myAfflictions + 1] = {name = resultingAffliction.name, 
					description = resultingAffliction.description,
					icon = affIcon}
					
				send("rendOdinCharacterClassHandler",
					"odinRendererCharacterAddAffliction",
					SELF.id,
					#state.myAfflictions, 
					resultingAffliction.name,
					resultingAffliction.description,
					affIcon)
					
			else
				state.myAfflictions[1] = {name = resultingAffliction.name, 
					description = resultingAffliction.description,
					icon = affIcon}
					
				send("rendOdinCharacterClassHandler",
					"odinRendererCharacterAddAffliction",
					SELF.id,
					#state.myAfflictions,
					resultingAffliction.name,
					resultingAffliction.description,
					affIcon)
			end			
		end

		local maxAfflictions = 1
		if SELF.tags["human"] or SELF.tags["fishperson"] then
			maxAfflictions = EntityDB.HumanStats.numAfflictionsToDie
		end
		if SELF.tags["bandit"] then
			maxAfflictions = EntityDB.Bandit.numAfflictionsToDie
		end

		-- TODO: WE SHOULD HAVE A UNIFIED SYSTEM FOR PULLING MAXAFFLICTIONS

		if #state.myAfflictions >= maxAfflictions and not SELF.tags.dead then
			-- check for death because there can be delayed attacks, and we don't want the character to "die" repeatedly.

			send(SELF, "deathBy", damagingObject, damageType)
			
			if damagingObject and
				damagingObject ~= SELF and
				not SELF.tags.selenian then -- TODO: clean up this process.
				
				send(damagingObject, "detectKilled", SELF, damageType)
			end
		else
			-- had affliction made, didn't die. Make icon happen.
			-- reset health, please.
			send(SELF, "emoteAffliction")
			state.AI.ints["health"] = state.AI.ints["healthMax"]
		end
	>>

	respond isDead()
	<<
		if SELF.tags["dead"] then
			return "isDeadResponse", true
		else
			return "isDeadResponse", false
		end
	>>

	respond hasAffliction(string name)
	<<
		for i, affliction in ipairs(state.myAfflictions) do
			if affliction.name == name then
				return "hasAfflictionResponse", true
			end
		end
		
		return "hasAfflictionResponse", false
	>>
	
	receive deathBy( gameObjectHandle damagingObject, string damageType )
	<<
		send(SELF,"ForceDropEverything")
	>>

	receive emoteAffliction()
	<<
		-- this is so we don't stack the affliction icon emote when taking a lot of damage in quick succession
		-- In the future, this might want to emote a particular affliction icon. 'Til then.
		
		if SELF.tags.human or
			SELF.tags.bandit or
			SELF.tags.fishperson then
			
			if state.AI.ints["emoteTimer"] > 4 then
				send("rendOdinCharacterClassHandler",
					"odinRendererCharacterExpression",
					state.renderHandle,
					"thought",
					"affliction",
					true )
				
				send(SELF, "resetEmoteTimer")
			end
		end
	>>
	
	-- CECOMMPATCH feature. Invisible fire for general use, or for entities that crash when the particles are attached
	receive invisFire()
	<<
		if SELF.tags["burning"] ~= true and
			not SELF.tags.dead and not
			SELF.tags.fire_immune then
			
			printl("CECOMMPATCH - invisible fire... oooOOooOooo spooktacular")
			
			-- game looks for a burning tag...
			SELF.tags["burning"] = true
			SELF.tags["burning_128"] = true
			
			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"SmallSplosion",
				state.AI.position.x,
				state.AI.position.y)
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererPlaySFXOnInteractive",
				state.renderHandle,
				"Fire Whoosh")
		end
	>>
	
	receive IgniteMessage()
	<<		
		-- CECOMMPATCH bugfix. Spores that die while on fire cause a crash, so do nothing if ignite attempted
		if SELF.tags["selenian_spore"] then
			send(SELF,"invisFire")
			return
		end
		
		-- use "Waist" for humanoids, "Root" as default otherwise.
		if SELF.tags["burning"] ~= true and
			not SELF.tags.dead and not
			SELF.tags.fire_immune then
			
			
			printl("ai_agent","holy cats, my name is " .. state.renderHandle .. " and I'm on fire!")

			if SELF.tags.citizen then --SELF.tags.human or SELF.tags.citizen then 
				send("rendCommandManager",
					"odinRendererTickerMessage",
					state.AI.name .. " has caught on fire!",
					"i_am_on_fire",
					"ui\\thoughtIcons.xml")
			end
			
			-- game looks for a burning tag...
			SELF.tags["burning"] = true
			SELF.tags["burning_128"] = true
			
			send("rendCommandManager",
				"odinRendererCreateParticleSystemMessage",
				"SmallSplosion",
				state.AI.position.x,
				state.AI.position.y)
			
			local attachmentPoint = "Root"
			
			if SELF.tags.human then
				attachmentPoint = "Waist"
			end
			
			send("rendOdinCharacterClassHandler",
				"odinRendererToggleParticlesOnJointMessage",
				state.renderHandle,
				attachmentPoint,
				"FirePerson", --"FireObject",
				true)
			
			send("rendInteractiveObjectClassHandler",
				"odinRendererPlaySFXOnInteractive",
				state.renderHandle,
				"Fire Whoosh")
		end		
	>>
	
	receive putOutFire()
	<<
		if SELF.tags.burning then
			printl("ai_agent","This part is over; " .. state.renderHandle .. " should no longer be on fire.")
			
			local attachmentPoint = "Root"
			if SELF.tags.human then
				attachmentPoint = "Waist"
			end
			
			send ("rendOdinCharacterClassHandler",
				"odinRendererToggleParticlesOnJointMessage",
				state.renderHandle,
				attachmentPoint,
				"FirePerson", --"FireObject",
				false)
			
			SELF.tags["burning"] = nil
			SELF.tags["burning_32"] = nil	
			SELF.tags["burning_64"] = nil
			SELF.tags["burning_128"] = nil
			SELF.tags["burning_256"] = nil
			SELF.tags["burning_512"] = nil
		end
	>>
>>