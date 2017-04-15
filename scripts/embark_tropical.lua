printl("SETTING TROPICAL BIOME")
local biomeName = "tropical"

send("gameSession", "setSessionString", "biome", biomeName)
send("gameSession", "setSessionBool", "biomeCold", false)
send("gameSession", "setSessionBool", "biomeTemperate", false)
send("gameSession", "setSessionBool", "biomeTropical", true)
send("gameSession", "setSessionBool", "biomeDesert", false)
send("gameSession", "setSessionInt", "colonyPopulation", 0)
send("gameSession", "setSessionInt", "militaryCount", 0)

send("rendCommandManager", "odinRendererSetLighting",126,132,183,33,37,82) -- Tropical Lighting

-- set up agriculture for this biome
-- (these will be read by farm.go to set allowed crops)

local cropTable = EntityDB.WorldStats.climateInfoPerBiome[ biomeName ].cropTable
for cropName, stats in pairs( cropTable ) do
     send("gameSession", "setSessionBool", "cropUnlocked=" .. cropName, cropTable[cropName].unlocked )
	send("gameSession", "setSessionInt", "cropGrowthModifier=" .. cropName, cropTable.growthModifier)
end

function spawnGameobject( x, y, objectType, objectTable )
	local createResults = query( "scriptManager", "scriptCreateGameObjectRequest", objectType, objectTable )
	local handle = createResults[1]
	if handle ~= nil then
		send(handle, "GameObjectPlace", x, y )
	end
end

-- need some accessible hunting at game start
local animals_to_spawn = { [1] = {["legacyString"]="Junglefowl"}, [2] = {["legacyString"]="Giant Beetle"} }
spawnGameobject( 250, 210, "herd", animals_to_spawn[ rand(1,#animals_to_spawn) ])
spawnGameobject( 250, 300, "herd", animals_to_spawn[ rand(1,#animals_to_spawn) ])

local spawnLocs = {
				[1] = {x=255,y=96, w=128, h=32 },
				[2] = {x=384,y=128, w=32, h=255},
				[3] = {x=255,y=384, w=128, h=32 },
			}

for i=1,9 do
	local spawnBox = spawnLocs[ rand(1,#spawnLocs) ]
	local spawnLoc = {x = rand(spawnBox.x, spawnBox.x + spawnBox.w),
				   y = rand(spawnBox.y, spawnBox.y + spawnBox.h) }
	
	spawnGameobject(spawnLoc.x,spawnLoc.y,"animal", {legacyString="Deathwurm"})
end
