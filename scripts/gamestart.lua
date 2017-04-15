--send("rendCommandManager", "odinRendererSetLighting", 255, 245, 235, 55, 60, 65) -- this is done in embarks now

--send("rendCommandManager", "SetSysBool", "desert_unlocked", true)

-- are traders on the map? No!
send("gameSession","setSessionBool","tradersOnMap",false) 

-- this is used to generate on-map ore veins but not used thereafter (yet)
local oremap = query("scriptManager",
				"scriptCreateGameObjectRequest",
				"oreMap", 
				{legacyString = "oremap"} )[1]

send("gameSession","setSessiongOH","oremap",oremap)

local analytics = query("scriptManager",
				"scriptCreateGameObjectRequest",
				"analytics", 
				{legacyString = "analytics"} )[1]

send("gameSession","setSessiongOH","analytics",analytics)

send("gameSession","setSessionInt","buildingCount",0)
send("gameSession","setSessionInt","disturbancePoints",0)

send("gameSession","setSessionInt","x_max", 255)
send("gameSession","setSessionInt","y_max", 255)

-- for overworld integration controls
-- this is an example and this whole block should be treated as debug
--[[local overworld_tag_list = {
		more_stahlmark = true,
		more_novorus = false,
		more_republique = false,
		bandits = true,
		more_bandits = true,
		obeliskians = false,
		more_obeliskians = false,
		fishpeople = true,
		more_fishpeople = true,
		mineral_poor = false,
		poor_farming = true, }

for k,v in pairs(overworld_tag_list) do
	send("gameSession","setSessionBool",k,v)
end]]

send("gameSession","setSessionBool","bandits",true)
send("gameSession","setSessionBool","fishpeople",true)
send("gameSession","setSessionBool","obeliskians",true)

send("gameSession","setSessionBool","more_stahlmark",false)
send("gameSession","setSessionBool","more_novorus",false)
send("gameSession","setSessionBool","more_republique",false)

--send("gameSession", "setSteamStatValue", "cannibalism", 50);
--send("gameSession", "setSteamAchievement", "cannibalismHasOccurred")
-- for endgame screen.
for i=1,20 do
	send("gameSession","setSessionString","endGameString" .. i ,"0")
end

--[[
for i=1,150 do
	send("gameSession", "incSteamStat", "stat_dodo_kill_count", 1)
end
--]]

send("gameSession","setSessionInt","militaryDeaths",0)
send("gameSession","setSessionInt","colonyDeaths",0)

send("gameSession","setSessionInt","highestPopulation",0) 			-- 1 DONE
send("gameSession","setSessionInt","totalImmigrantCount",0) 		-- 2 DONE
send("gameSession","setSessionInt","highestBuildingCount",0) 		-- 3 DONE
send("gameSession","setSessionInt","totalGoodsProduced",0) 			-- 4 DONE
send("gameSession","setSessionInt","totalAgriculturalOutput",0) 		-- 5 DONE
send("gameSession","setSessionInt","treesChoppedCount", 0) 			-- 6 DONE
send("gameSession","setSessionInt","actsOfCannibalism", 0) 			-- 7 DONE
send("gameSession","setSessionInt","highestCultistPopulation", 0) 	-- 8 DONE
send("gameSession","setSessionInt","majorEldritchEvents", 0) 		-- 9 DONE
send("gameSession","setSessionInt","eldritchArtifactsFound", 0) 		-- 10 DONE

send("gameSession", "setSessionInt", "tipTopCaviarCount", 0)
send("gameSession", "setSessionBool", "cannibalismHasOccurred", false)
send("gameSession", "setSessionInt", "oreNodesMinedCount", 0)

send("gameSession", "setSessionInt", "workplaceCount", 0)
local overseerCount = 7 -- for now this is what it starts at.
local workplaceCap = overseerCount - 1
send("gameSession", "setSessionInt", "workplaceCap", workplaceCap)


local x_max = query("gameSession","getSessionInt","x_max")[1]
local y_max = query("gameSession","getSessionInt","y_max")[1]

send("rendCommandManager",
	"odinRendererTickerMessage",
	"Beginning Day 1",
	"sun",
	"ui\\thoughtIcons.xml")

-- do factions setup.
for k,v in pairs(EntitiesByType["faction"]) do
	local create_results = query( "scriptManager",
							"scriptCreateGameObjectRequest",
							"faction", 
							{legacyString = k} )
	
	local handle = create_results[1]
end

-- do fishpeople faction object.
local create_results = query( "scriptManager",
						"scriptCreateGameObjectRequest",
						"faction", 
						{legacyString = "FishInfo" } )

--send("gameSession", "setSessionInt", "upkeep_trunks1", 0)
--send("gameSession", "setSessionInt", "upkeep_trunks2", 0)
--send("gameSession", "setSessionInt", "upkeep_trunks3", 0)
--send("gameSession", "setSessionInt", "upkeep_trunks4", 0)


-- science stats to be changed by Science! Research!
	--[[
	industry modifiers: consider them like skills that apply to the entire colony
	int should be treated as a percentage, so 100 is baseline.
	to INCREASE efficiency SUBTRACT from the modifier. 95 = 95% of time to do a job.
	Tread lightly here because they modify skill times, which can go very low.
	]]
	
	send("gameSession", "setSessionInt", "agricultureTechModifier", 100)
	send("gameSession", "setSessionInt", "miningTechModifier", 100)
	send("gameSession", "setSessionInt", "carpentryTechModifier", 100)
	send("gameSession", "setSessionInt", "cookingTechModifier", 100)
	send("gameSession", "setSessionInt", "chemistryTechModifier", 100)
	send("gameSession", "setSessionInt", "stoneworkingTechModifier", 100)
	send("gameSession", "setSessionInt", "smeltingTechModifier", 100)
	send("gameSession", "setSessionInt", "metalworkingTechModifier", 100)
	send("gameSession", "setSessionInt", "researchTechModifier", 100)
	send("gameSession", "setSessionInt", "diplomacyTechModifier", 100)
	send("gameSession", "setSessionInt", "militaryTrainingTechModifier", 100)
	send("gameSession", "setSessionInt", "militaryDamageTechBonus", 0)
	send("gameSession", "setSessionInt", "militaryDefenseTechBonus", 0)
	send("gameSession", "setSessionInt", "gatheringSpeedModifier", 100)
	send("gameSession", "setSessionInt", "militaryHealthBonus", 0)
	send("gameSession", "setSessionInt", "civilianHealthBonus", 0)
	send("gameSession", "setSessionInt", "militaryReloadModifier", 100)
	
	--send("gameSession", "setSessionInt", "medicalTechModifier", 100) -- not done
	--send("gameSession", "setSessionInt", "constructionTechModifier", 100) -- not done, might who cares
	
	-- science unlocks
	--[[ format for crop modifiers/unlocks:
	
		send("gameSession", "setSessionBool", "cropUnlocked=" .. cropName, BOOL )
		send("gameSession", "setSessionInt", "cropGrowthModifier=" .. cropName, INT )
		
		note: these are set up automatically per-embark type using info from
			WorldStats.climeInfoPerBiome[biomeName].cropTable
		]]
		
	-- Science study counters.
	
	send("gameSession","setSessionInt","naturalistMineNodesStudied",0)
	
-- day/night cycle
	send("gameSession","setSessionString","eventLighting", "default") --used to toggle on event lighting if set to something other than default.
	send("gameSession","setSessionString","lastLightingType", "default") --used to tell the lighting transition system what the current lighting type is.
	send("gameSession", "setSessionBool", "refreshLighting", false) --If true, frameupdate will instantly reload the lighting on the current zone. Use for activating event lighting.
	send("gameSession", "setSessionInt", "dayCount", 1)
	send("gameSession", "setSessionInt",  "dayNightCounter", 0)
	send("gameSession", "setSessionBool", "transitioning", false)
	send("gameSession", "setSessionInt",  "transitionCounter", 0)
	send("gameSession", "setSessionInt", "currentShift", 1)
	send("gameSession", "setSessionString", "lightingTransitionFrom", "sunrise")
	send("gameSession", "setSessionString", "lightingTransitionTo", "sunrise")
	send("gameSession", "setSessionBool", "isNight", false)
	send("gameSession", "setSessionBool", "isDay", true)
	send("gameSession", "setSessionBool", "isDawn", false)
	send("gameSession", "setSessionBool", "isDusk", false)
	send("gameSession", "setSessionBool", "isNoon", false)

	send("gameSession", "setSessionBool", "shiftTimer", 0)


-- There was a prestige here. It's gone now. 

send("gameSession", "setSessionInt", "airdropX", 127)
send("gameSession", "setSessionInt", "airdropY", 127)
send("gameSession", "setSessionBool", "airshipMastBuilt", false)
send("gameSession", "setSessionBool", "airdropOverride", false)

-- foreign relations
	-- -100 to 100 scale relations
	-- bools for allied & war
	
	send("gameSession", "setSessionInt", "NovorusRelations", 0)
	send("gameSession", "setSessionInt", "NovorusNeutralHostile", -33)
	send("gameSession", "setSessionInt", "NovorusNeutralFriendly", 33)
	send("gameSession", "setSessionInt", "NovorusRelationsMin", -100)
	send("gameSession", "setSessionInt", "NovorusRelationsMax", 100)
	send("gameSession", "setSessionBool", "novorus_max", false)
	send("gameSession", "setSessionInt", "NovorusLoggingBan", 0)

	send("gameSession", "setSessionInt", "RepubliqueRelations", 0)
	send("gameSession", "setSessionInt", "RepubliqueNeutralHostile", -33)
	send("gameSession", "setSessionInt", "RepubliqueNeutralFriendly", 33)
	send("gameSession", "setSessionInt", "RepubliqueRelationsMin", -100)
	send("gameSession", "setSessionInt", "RepubliqueRelationsMax", 100)
	send("gameSession", "setSessionBool", "republique_max", false)
	send("gameSession", "setSessionInt", "RepubliqueCrafting", 0)
	send("gameSession", "setSessionBool", "AllowUraniumCrafting", false)
	
	send("gameSession", "setSessionInt", "StahlmarkRelations", 0)
	send("gameSession", "setSessionInt", "StahlmarkNeutralHostile", -33)
	send("gameSession", "setSessionInt", "StahlmarkNeutralFriendly", 33)
	send("gameSession", "setSessionInt", "StahlmarkRelationsMin", -100)
	send("gameSession", "setSessionInt", "StahlmarkRelationsMax", 100)
	send("gameSession", "setSessionBool", "stahlmark_max", false)
	send("gameSession", "setSessionInt", "StahlmarkCrafting", 0)

	send("gameSession", "setSessionInt", "EmpireRelations", 0)
	send("gameSession", "setSessionInt", "EmpireRelationsMin", -100)
	send("gameSession", "setSessionInt", "EmpireRelationsMax", 100)
	send("gameSession", "setSessionBool", "invasionInProgress", false)

	--send("gameSession","setSessiongOH", "cult_shrine", nil)
	
	-- TEST: at gamestart, make one hostile, one neutral, one allied.
	
	local nations = { "Stahlmark", "Novorus", "Republique" }
	local r = rand(1,3)
	local hostileNation = nations[r]
	table.remove(nations, r)
	r = rand(1,2)
	local alliedNation = nations[r]
	table.remove(nations, r)
	local neutralNation = nations[1]
	
	printl("events", "making hostile nation out of: " .. hostileNation)
	
	send("gameSession", "setSessionInt", hostileNation .. "Relations", -75)
	send( query("gameSession","getSessiongOH", hostileNation)[1], "makeHostile")
	
	printl("events", "making potentially allied nation out of: " .. alliedNation)
	--send("gameSession", "setSessionInt", alliedNation .. "Relations", 75)
	--send( query("gameSession","getSessiongOH", alliedNation)[1], "makeFriendly")
	send( query("gameSession","getSessiongOH", alliedNation)[1], "makeNeutral")
	
	send( query("gameSession","getSessiongOH", neutralNation)[1], "makeNeutral")

	send("gameSession", "setSessionString", "defaultFriendly", alliedNation)
	send("gameSession", "setSessionString", "defaultHostile", hostileNation)
	send("gameSession", "setSessionString", "defaultNeutral", neutralNation)

-- just #fishpeoplething

	-- this is used to set delays between fishpeople reaction/policy events.
	send("gameSession", "setSessionBool", "fishpeopleEventActive", false)

	send("gameSession", "setSessionBool", "fishpeopleFirstContact", false)
	send("gameSession","setSessionBool","fishpeopleShotOnSight",false)

	send("gameSession", "setSessionBool", "fishpeoplePolicyHostile", false)
	send("gameSession", "setSessionBool", "fishpeoplePolicyDenial", false)
	send("gameSession", "setSessionBool", "fishpeoplePolicyFriendly", false)
	
	send("gameSession", "setSessionInt", "fishpeoplePolicyLastSetDay", 0)

	send("gameSession", "setSessionInt", "fishpeopleConversations", 0)
	send("gameSession", "setSessionInt", "fishpeopleAnger", 0)
	send("gameSession", "setSessionInt", "fishpeopleDeaths", 0)
	send("gameSession", "setSessionInt", "lastFishpeopleRaidDay", 0)
	send("gameSession", "setSessionInt", "fishpeoplePopulation", 0) -- on map
	
	send("gameSession","setSessionBool","fishpeopleVandalismPolicyEventActive",false)
	send("gameSession","setSessionBool","fishpeopleVandalismExecution",false)
	send("gameSession", "setSessionBool","fishpeopleVandalismDiscouraged",false)
	
	send("gameSession", "setSessionBool", "fishpeopleButcherHumanPolicySet", false)
	send("gameSession", "setSessionBool", "fishpeopleButcherHumanPolicyDeath", false)
	send("gameSession", "setSessionBool", "fishpeopleButcherHumanPolicyBeatings", false)
	
	send("gameSession", "setSessionBool", "fishpeopleHasslePolicySet", false)
	send("gameSession", "setSessionBool", "fishpeopleHassleAll", true)
	send("gameSession", "setSessionBool", "fishpeopleHassleTroublemakers", true)
	
-- Obeliskians
	send("gameSession", "setSessionInt", "obeliskiansOnMap", 0)
	send("gameSession", "setSessionBool", "obeliskiansFirstContact", false)
	send("gameSession", "setSessionInt", "obeliskianDeaths", 0)
	
-- Selenians
	send("gameSession", "setSessionInt", "seleniansOnMap", 0)
	
-- bandits
	send("gameSession", "setSessionInt", "banditRegionPool", 100)
	send("gameSession", "setSessionInt", "banditsLastSpawnDay",0)
	send("gameSession", "setSessionInt", "banditsOnMap", 0)
	send("gameSession", "setSessionInt", "banditDeaths", 0)

	send("gameSession", "setSessionInt", "timesCapitulatedToBandits", 0)
	send("gameSession", "setSessionBool", "banditTruceEventFired", false)
	
	send("gameSession", "setSessionBool", "banditPlunderDone", false)
	send("gameSession", "setSessionBool", "banditPlunderDenied", false)
	send("gameSession", "setSessionInt", "plunderAcceptanceDay", 0)
	
	-- new bandit control variables. DGB to clean up everything above.
	-- Bandits will start at fairly hostile.
	send("gameSession", "setSessionBool", "BanditsHostile", true)
	send("gameSession", "setSessionBool", "BanditsNeutral", false)
	send("gameSession", "setSessionBool", "BanditsFriendly", false)
	send("gameSession", "setSessionInt", "BanditsRelations", -50)
	send("gameSession", "setSessionInt", "BanditsNeutralHostile", -33)
	send("gameSession", "setSessionInt", "BanditsNeutralFriendly", 33)
	send("gameSession", "setSessionInt", "BanditsRelationsMin", -100)
	send("gameSession", "setSessionInt", "BanditsRelationsMax", 100)

	send("rendCommandManager","SetFactionPolicyString", "Bandits", "Hostile")

	
	send("gameSession","setSessionInt","BanditsLastSpawnDay", 0)
	send("gameSession","setSessionInt","BanditsLastPlunderDay", 0)
	send("gameSession","setSessionBool","BanditsFirstContact", false)
	send("gameSession", "setSessionBool", "BanditsForcedHostility", false)
	
	-- burial.
	send("gameSession","setSessionBool","BanditsBuryCorpses",false)
	send("gameSession","setSessionBool","BanditsDumpCorpses",false)
	send("gameSession","setSessionBool","BanditsBuryCorpsesPolicySet",false)
	send("gameSession","setSessionInt","BanditsCorpsePolicyDay",0)
	
	-- for testing
	send("gameSession", "setSessionBool", "testingPlaceholder", false) 
 	
	
send("gameSession","setSessionInt","steamKnightsActive",0)
	
-- random horrors

send("gameSession", "setSessionInt", "spectreLastReportDay", 0)
send("gameSession", "setSessionInt", "lastPopulationEventDay", -1)

--event arc stuff
send("gameSession", "setSessionBool", "trackClearables", false)

-- LC/MC/HC population counts are initialized in the loadout scripts because those are executed first & contain the pop spawn scripts.
-- The session ints, if you're here looking for them, are named lowerClassPopulation, middleClassPopulation.
send("gameSession", "setSessionInt", "LcPopulationCap", 112)
send("gameSession", "setSessionInt", "McPopulationCap", 38)
--send("gameSession", "setSessionInt", "LcPopulationAllowed", 2) these are done in loadout now
--send("gameSession", "setSessionInt", "McPopulationAllowed", 7)
--send("gameSession", "setSessionInt", "totalPopulationAllowed", 7)
send("gameSession", "setSessionBool", "immigrationStatus", true)
send("gameSession", "setSessionBool", "mcImmigrationStatus", true)
send("gameSession", "setSessionBool", "crisisStatus", false)
send("gameSession", "setSessionInt", "immigrationTimes", 0) --this is for the new immigration sliding timer

send("gameSession", "setSessionInt", "tempCharacterPopulation", 0)
send("gameSession", "setSessionInt", "permCitizensDead", 0)
send("gameSession", "setSessionInt", "tempCitizensDead", 0)
send("gameSession", "setSessionInt", "deadOverseers", 0)
send("gameSession", "setSessionInt", "replacedOverseers", 0)

send("gameSession", "setSessionInt", "militaryCount", 0) -- how many soldiers are in the colony?


-- cult vars -- probably deprecated well before 52B
send("gameSession", "setSessionInt", "cultPower", 0)
send("gameSession", "setSessionInt", "cultEvents", 0)
send("gameSession", "setSessionBool", "cultShrineDecisionMade", false)
send("gameSession", "setSessionBool", "cultShrineDecisionInProgress", false)
send("gameSession", "setSessionBool", "cultShrineDecisionWaffle", false)
send("gameSession", "setSessionBool", "tolerateCultShrines", false)
send("gameSession", "setSessionBool", "persecuteCultShrines", false)
send("gameSession", "setSessionInt", "numberOfCultShrines", 0)
send("gameSession", "setSessionInt", "lastMinistryCultInvestigation", 0)
send("gameSession", "setSessionBool", "creepyCultEventSeen", false)

send("gameSession", "setSessionInt", "starvingCount", 0)
send("gameSession", "setSessionBool", "starvationBailout", false)

--[[send("gameSession", "setSessionInt", "dodoCount", 0)
send("gameSession", "setSessionBool", "dodoExtinction", false)
send("gameSession", "setSessionInt", "aurochsCount", 0)
send("gameSession", "setSessionBool", "aurochsExtinction", false)--]]


-- tutorial 
	--[[ old tutorial stuff
	send("gameSession", "setSessionBool", "tutorialCurrentlyActive", false)
	send("gameSession", "setSessionBool", "enableContextualTutorials", false)
	send("gameSession", "setSessionBool", "farmTutorialDone", false)
	send("gameSession", "setSessionBool", "jobFilterTutorialDone", false)
	send("gameSession", "setSessionBool", "standingorderTutorialDone", false)
	send("gameSession", "setSessionBool", "moduleTutorialDone", false)
	send("gameSession", "setSessionBool", "stockpileTutorialDone", false)
	send("gameSession", "setSessionBool", "refiningTutorialDone", false)
	send("gameSession", "setSessionBool", "tutorialCultsDone", false)
	send("gameSession", "setSessionBool", "workshopTutorialDone", false)
	send("gameSession", "setSessionBool", "qualityTutorialDone", false)
	send("gameSession", "setSessionBool", "starterEventFired", false)
	]]
	send("gameSession", "setSessionInt", "jumpToTutorial", 0)

	send("gameSession", "setSessionBool", "caseTutorialActive", false)
	send("gameSession", "setSessionBool", "bedTutorialDone", false)
	
	send("gameSession", "setSessionBool", "barracksWarned1", false)
	send("gameSession", "setSessionBool", "barracksWarned2", false)
	send("gameSession", "setSessionBool", "townhallPromptWarned", false)
	
-- other events
send("gameSession", "setSessionBool", "30charEventOver", false)
send("gameSession", "setSessionBool", "50charEventOver", false)
send("gameSession", "setSessionBool", "70charEventOver", false)
send("gameSession", "setSessionBool", "100charEventOver", false)
send("gameSession", "setSessionBool", "ncoReplacementEvent", false)
send("gameSession", "setSessionInt", "vicarCount", 0)
send("gameSession", "setSessionBool", "vicarReplacementEvent", false)
send("gameSession", "setSessionBool", "chapelBuilt", false)
send("gameSession", "setSessionBool", "triggeredChapel", false)
send("gameSession", "setSessionBool", "macroscopeBuilt", false)
send("gameSession", "setSessionBool", "Immigration_MC_FirstTime", false)
send("gameSession", "setSessionBool", "MC_Bailout_Done", false)
send("gameSession", "setSessionBool", "noBedsWarningDone", false)
send("gameSession", "setSessionBool", "blockWeatherMemories", false)
--send("gameSession", "setSessionBool", "Immigration_MC_GetNCO", false)

-- ECONOMY
send("gameSession", "setSessionInt", "workshopsBuilt", 0)
send("gameSession", "setSessionInt", "barracksBuilt", 0)
send("gameSession", "setSessionInt", "townHallsBuilt", 0)
send("gameSession", "setSessionInt", "tier0produced", 0) 
send("gameSession", "setSessionInt", "tier1produced", 0)
send("gameSession", "setSessionInt", "tier2produced", 0)
send("gameSession", "setSessionInt", "tier3produced", 0)
send("gameSession", "setSessionInt", "tier4produced", 0)
send("gameSession", "setSessionInt", "tier5produced", 0)
send("gameSession", "setSessionInt", "tier6produced", 0)
send("gameSession", "setSessionInt", "tier7produced", 0)
send("gameSession", "setSessionInt", "nextEconomyGoal", 1)
send("gameSession", "setSessionInt", "foodTier0produced", 0)
send("gameSession", "setSessionInt", "foodTier1produced", 0)
send("gameSession", "setSessionInt", "UCHousesProduced", 0)


--ACHIEVEMENT STORAGE

     --Buildings the player can build:
     send("gameSession","setSessionBool","builtTraining Academy", "false")
     send("gameSession","setSessionBool","builtBarracks", "false")
     send("gameSession","setSessionBool","builtTrade Office", "false")
     send("gameSession","setSessionBool","builtForeign Office", "false")
     send("gameSession","setSessionBool","builtNaturalist's Office", "false")
     send("gameSession","setSessionBool","builtLaboratory", "false")
     send("gameSession","setSessionBool","builtChemical Works", "false")
     send("gameSession","setSessionBool","builtKitchen", "false")
     send("gameSession","setSessionBool","builtPublic House", "false")
     send("gameSession","setSessionBool","builtUpper Class House", "false")
     send("gameSession","setSessionBool","builtMiddle Class House", "false")
     send("gameSession","setSessionBool","builtLower Class House", "false")
     send("gameSession","setSessionBool","builtMine", "false")
     send("gameSession","setSessionBool","builtCeramics Workshop", "false")
     send("gameSession","setSessionBool","builtMetalworks", "false")
     send("gameSession","setSessionBool","builtCarpentry Workshop", "false")
	
	send("gameSession","setSessionInt","foodFarmCount", 0)
     


function spawnGameobject( x, y, objectType, objectTable )

	if x > x_max - 4 then x = x_max - 4 end
	if x < 4 then x = 4 end
	
	if y > y_max - 4 then y = y_max - 4 end
	if y < 4 then y = 4 end

	local createResults = query("scriptManager",
						   "scriptCreateGameObjectRequest",
						   objectType,
						   objectTable )
	
	local handle = createResults[1]
	
	send(handle,
		"GameObjectPlace",
		x,
		y  )
	
	--local radius = 2
	--if handle ~= nil then
		
	--	local added = false
		--[[while not added do
				
			local new_location = gameGridPosition:new()
			new_location.x = x + rand(-radius, radius)
			new_location.y = y + rand(-radius, radius)
		
			printl("DAVID", "trying: " .. tostring(new_location.x) .. " / " .. tostring(new_location.y) )
			local gSD_results = query("gameSpatialDictionary",
								 "gridCanAddObjectTo", 
		                          	  handle, 
		                          	  new_location )
			
			local iswater = query( "gameSpatialDictionary",
								"gridHasSpatialTag",
								new_location,
								"water" )[1]

			if gSD_results[1] then
				if not iswater then
					send( handle,
						"GameObjectPlace",
						new_location.x,
						new_location.y  )
					
					added = true
				end
			end
			radius = radius + 1
		end]]
	--end

	return handle
end

local startX = query("gameSession", "getSessionInt", "startX")[1]
local startY = query("gameSession", "getSessionInt", "startY")[1]
local newSpawnPoint = gameGridPosition:new()

newSpawnPoint.x = startX
newSpawnPoint.y = startY

send("gameSpatialDictionary", "gridSetPlayerSpawnPoint", newSpawnPoint, 10  )

-- artifacts
spawnGameobject( 	rand(32, x_max - 16),
				rand(32,y_max - 16),
				"clearable",
				{ legacyString="A Mundane Pile Of Dirt" } )

spawnGameobject( 	rand(32,x_max - 16),
				rand(32,y_max - 16),
				"clearable",
				{ legacyString="A Mundane Pile Of Dirt" } )

-- Obeliskians! -- clean up these positions.
if query("gameSession","getSessionBool","obeliskians")[1] == true then 
	spawnGameobject( 	rand(16, math.floor(x_max *0.4) - 8),
					rand(16,y_max - 16),
					"objectcluster",
					{ legacyString="Obeliskian Cluster" } )
	
	spawnGameobject( 	rand(16,x_max - 16),
					rand(16, math.floor(y_max *0.4) - 8),
					"objectcluster",
					{ legacyString="Obeliskian Cluster" } )
	
	if query("gameSession","getSessionBool","more_obeliskians")[1] == true then 
		spawnGameobject( 	rand(16, math.floor(x_max *0.4) - 8),
						rand(16,y_max - 16),
						"objectcluster",
						{ legacyString="Obeliskian Cluster" } )
		
		spawnGameobject( 	rand(16,x_max - 16),
						rand(16, math.floor(y_max *0.4) - 8),
						"objectcluster",
						{ legacyString="Obeliskian Cluster" } )
	end
end

-- set up event directors & event director points
send("gameSession", "setSessionInt", "arcPointsPoolEldritch", 0)
send("gameSession", "setSessionInt", "arcPointsPoolMundane", 0)

--send("gameSession", "setSessionInt", "arcPointsGainPerDayEldritch", EntityDB.WorldStats.arcPointDripDefaultEldritch ) 
--send("gameSession", "setSessionInt", "arcPointsGainPerDayMundane", EntityDB.WorldStats.arcPointDripDefaultMundane ) 

local eldritch_director = query("scriptManager",
				"scriptCreateGameObjectRequest",
				"event_director",
				{ name = "eldritch" } )[1]

local mundane_director = query("scriptManager",
				"scriptCreateGameObjectRequest",
				"event_director",
				{ name = "mundane" } )[1]

local other_director = query("scriptManager",
				"scriptCreateGameObjectRequest",
				"event_director",
				{ name = "other" } )[1]

-- clay & stone to be nice
spawnGameobject( startX + rand(-24,24),  startY + rand(-20,20), "objectcluster", { legacyString="Clay Nodes" } )
spawnGameobject( startX + rand(-28,28),  startY + rand(-28,28), "objectcluster", { legacyString="Clay Nodes" } )
spawnGameobject( startX + rand(-36,36),  startY + rand(-36,36), "objectcluster", { legacyString="Clay Nodes" } )

spawnGameobject( startX + rand(-20,20), startY + rand(-20,20), "objectcluster", { legacyString="Stone Boulders" } )
spawnGameobject( startX + rand(-28,28), startY + rand(-28,28), "objectcluster", { legacyString="Stone Boulders" } )

send("gameSpatialDictionary", "gridExploreFogOfWar", startX, startY, 45)

-- set up airdrop controller
local result = query("gameObjectManager", "gameObjectCollectionRequest", "airdropMasts")
local collection = result[1]
collection = {}

-- set up Steam Knight collection
local sk_collection = query("gameObjectManager", "gameObjectCollectionRequest", "steamKnights")[1]
sk_collection = {}

local handle2 = query( "scriptManager",
		"scriptCreateGameObjectRequest",
		"event_director",
		{name="tutorial"} )[1]

send(handle2,"startNewEventArc","Tutorial")

send("gameSession","setSessionBool","horror_policy_study", false)
send("gameSession","setSessionBool","horror_policy_harvest", false)
send("gameSession","setSessionBool","horror_policy_dump", true)

printl("events", "dominantFaction = " .. query("gameSession","getSessionString", "dominantFaction")[1] )
