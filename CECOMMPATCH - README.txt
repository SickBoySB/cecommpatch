
   _____ ______ _____ ____  __  __ __  __ _____     _______ _____ _    _ 
  / ____|  ____/ ____/ __ \|  \/  |  \/  |  __ \ /\|__   __/ ____| |  | |
 | |    | |__ | |   | |  | | \  / | \  / | |__) /  \  | | | |    | |__| |
 | |    |  __|| |   | |  | | |\/| | |\/| |  ___/ /\ \ | | | |    |  __  |
 | |____| |___| |___| |__| | |  | | |  | | |  / ____ \| | | |____| |  | |
  \_____|______\_____\____/|_|  |_|_|  |_|_| /_/    \_\_|  \_____|_|  |_|
  
            Clockwork Empires Community Patch - Version: 0.8.1
     Clockwork Empires, all code, assets, etc. are (c) Gaslamp Games

Welcome to the Clockwork Empires Community Patch (CECOMMPATCH)! This project is an effort to continue development of the game Clockwork Empires that was made by Gaslamp Games. All changes are done to the exposed game files, *not* the engine itself. Decompiling the executable to gain access to those parts of the code is not part of this project, so engine-related limitations exist. Fortunately, Gaslamp Games made their engine in such a way that the vast majority of the game is moddable so it's rare to run into hardcoded issues.

The first priority of CECOMMPATCH is to fix any outstanding bugs that the game has. Great progress has been made so far, and the list of bugs is ever decreasing. After the majority are fixed, the next step will be to implement features that were removed/unfinished. A huge amount of unused 3D models, events, and gameplay features are still in the game files, and we plan to utilize as much as possible.

After that stage, completely new content will be considered. The intent of this project is not to create a "total conversion", so all custom content will be themed in a way that fits with the established style.


  _____ _   _  _____ _______       _      _      
 |_   _| \ | |/ ____|__   __|/\   | |    | |     
   | | |  \| | (___    | |  /  \  | |    | |     
   | | | . ` |\___ \   | | / /\ \ | |    | |     
  _| |_| |\  |____) |  | |/ ____ \| |____| |____ 
 |_____|_| \_|_____/   |_/_/    \_\______|______|

            INSTALL

===============================================================
THIS IS FOR THE "1.0D - Experimental" VERSION OF THE GAME ONLY!
===============================================================

1. BUY Clockwork Empires from Steam! http://store.steampowered.com/app/224740/

2. Download the "Experimental Branch" (1.0D) version of Clockwork Empires from Steam
	1a. Install the game normally
	1b. Right click the game in your Steam Library
	1c. Click "Properties"
	1d. Click the "BETAS" tab
	1e. Choose "experimental-" from the dropdown
	1f. Close the pop-up and allow the game to update
	1g. The game should now be called "Clockwork Empires [experimental]" in your library

3. Download the latest CECOMMPATCH from https://github.com/SickBoySB/cecommpatch/archive/v1.zip

4. Open zip file, and open the "cecommpatch-1" folder inside

5. Unzip the contents into your game folder (which should be located wherever Steam installed it, typically "C:\Program Files (x86)\Steam\steamapps\common\Clockwork Empires")

6. When prompted to overwrite the files, say "YES"

7. You're done! It is highly recommended to use a new save with this patch. While initial versions should work okay with other 1.0D saves, as new content is added there will likely be problems.


To uninstall this patch all you need to do is "verify integrity of game files" on Steam. The original files will be redownloaded and the mod will be overwritten.


  _  ___   _  ______          ___   _   _____  _____ _____ _    _ ______  _____ 
 | |/ / \ | |/ __ \ \        / / \ | | |_   _|/ ____/ ____| |  | |  ____|/ ____|
 | ' /|  \| | |  | \ \  /\  / /|  \| |   | | | (___| (___ | |  | | |__  | (___  
 |  < | . ` | |  | |\ \/  \/ / | . ` |   | |  \___ \\___ \| |  | |  __|  \___ \ 
 | . \| |\  | |__| | \  /\  /  | |\  |  _| |_ ____) |___) | |__| | |____ ____) |
 |_|\_\_| \_|\____/   \/  \/   |_| \_| |_____|_____/_____/ \____/|______|_____/ 

            KNOWN ISSUES

- Spores on fire don't show the fire effect. This is the only way to fix the related crash.

- Only ~5 alerts are shown without scrolling. This is a tradeoff to fix alerts pushing the buttons off the screen.

- When swapping from one save to another the alerts menu disappears until a new alert occurs (or an old one goes away). This ONLY happens when swapping between saves, and a full restart of the application will not have this problem. This is a tradeoff to fix alerts pushing the buttons off the screen.

- The "Bottoms Up" and "Who You Gonna Call" achievements will not visually progress in Steam, but the achievement DOES work as intended now. The stat/progress portion is impossible to fix because the variable is unknown and all educated guesses didn't work. Since the stat wasn't being tracked internally before, only new booze served/spectres banished will count towards the achievements. Both stats are save-game specific, so it is NOT an overall amount.

- Blueprint requirements show the log icon instead of the new timber icon. This is because of how that part of the UI functions. The icon cannot be changed without breaking the requirements calculations.

- Pause and building toggling buttons on the top menu don't have the fancy hover effect the other menus have. This is unavoidable because of how those work on the hardcoded side of things. This is the same reason the factions and jobs buttons don't have the little arrow when they're selected.

  _      _____ _   _ _  __ _____ 
 | |    |_   _| \ | | |/ // ____|
 | |      | | |  \| | ' /| (___  
 | |      | | | . ` |  <  \___ \ 
 | |____ _| |_| |\  | . \ ____) |
 |______|_____|_| \_|_|\_\_____/ 

            LINKS

Steam Group: http://steamcommunity.com/groups/cecommpatch
The Steam Group is the ideal place to keep up to date with the project, as well as discuss changes and report bugs. It's a public group, so no invite is necessary. Please adhere to the rules set in place though, as this is *not* the place to bash Gaslamp Games or troll.

Official Thread: https://community.gaslampgames.com/threads/community-patch.20992/
This is the official thread created on the Gaslamp Games forum. There is no telling how long they will stay online, so it's probably better to stay in touch via Steam.

GitHub: https://github.com/SickBoySB/cecommpatch/
Feel free to fork the code, contribute your own changes, make your own patch, etc.


   _____ _    _          _   _  _____ ______ _      ____   _____ 
  / ____| |  | |   /\   | \ | |/ ____|  ____| |    / __ \ / ____|
 | |    | |__| |  /  \  |  \| | |  __| |__  | |   | |  | | |  __ 
 | |    |  __  | / /\ \ | . ` | | |_ |  __| | |   | |  | | | |_ |
 | |____| |  | |/ ____ \| |\  | |__| | |____| |___| |__| | |__| |
  \_____|_|  |_/_/    \_\_| \_|\_____|______|______\____/ \_____|

            CHANGELOG
         __     ____     __   
        /  \   / _  \   /  \  
       (  0 )_ ) _  ( _(_/ /    0.8.1
        \__/(_)\____/(_)(__)  

FIX: UI - Laboratory projects state "Already Completed" after swapping to another category
FIX: UI - Resetting laboratory projects in a non-selected group doesn't update until re-selected
CHANGE: UI - "Jobs" button changed to the fancy version
CHANGE: UI - "Factions" button change to the fancy version
CHANGE: UI - Population/workplace cap top menu item has been spruced up
CHANGE: UI - Disturbance points top menu item has been spruced up
CHANGE: UI - Provisions top menu item has been spruced up
CHANGE: UI - Toggle building view button changed to a fancy version
CHANGE: UI - Pause button changed to a fancy version
CHANGE: UI - Time of day image moved, ugly (pointless) border removed

         __     ____     __  
        /  \   / _  \   /  \ 
       (  0 )_ ) _  ( _(  0 )   0.8.0
        \__/(_)\____/(_)\__/ 

!!!MAJOR CHANGE!!! - Male colonists now have a chance to wear (intelligently) randomized hats!
FIX: UI - Blueprint requirements for timber are showing in red when there is amble timber
FIX: UI - "Ornate Table and Chair Set" is not showing up in the module placement menus
FIX: ICON - "MK-1 Steam Knight" uses prototype version's icon
CHANGE: UI - "All Decor" filter reimplemented
CHANGE: UI - "All Modules" filter reimplemented
CHANGE: UI - "All Doors/Windows" filter added
CHANGE: UI - "All Furniture" filter added
CHANGE: TAG - Farming added to Lacquer
CHANGE: TAG - Farming added to Bamboo
CHANGE: TAG - Icon for "Used in Workshops" changed to general workshop icon

         __    ____   __
        /  \  (__  ) /  \
       (  0 )_  / /_(  0 )   0.7.0   (previous changelogs combined)
        \__/(_)(_/(_)\__/

!!!MAJOR CHANGE!!! - Simplification and standardization of item names. Nearly every item is impacted!
!!!MAJOR CHANGE!!! - Tooltip tags have had a total overhaul in how they are used. EVERY item/recipe tooltip impacted!
FIX: CRASH - Selenian Spores killed while on fire. See "known issues" for info on the tradeoff
FIX: ACHIEVEMENT - "Bottoms Up" not triggering. See "known issues" for info on the tradeoff
FIX: ACHIEVEMENT - "Who You Gonna Call" not triggering. See "known issues" for info on the tradeoff
FIX: ACHIEVEMENT - "Tidy Estates" not triggering
FIX: UI - Icons pushed off the screen when there are too many alerts. See "known issues" for info on the tradeoff
FIX: UI - Icons pushed off the screen at low resolutions. See "known issues" for info on the tradeoff
FIX: UI - Decor/Module menu cutoff at the bottom. See "known issues" for info on the tradeoff
FIX: UI - Main menu cutoff at low resolutions
FIX: ICON - "MK-1 Machine Gun" wrong
FIX: ICON - "MK-1 Steam Knight Chassis" wrong
FIX: ICON - "MK-1 Steam Knight Oculars" wrong
FIX: ICON - "Ammo Autoloader" wrong
FIX: ICON - Gray "Mine Shorings" missing completely
FIX: ICON - Gray "MK-1 Steam Knight Oculars" missing completely
FIX: ICON - Gray "MK-1 Steam Knight Chassis" missing completely
FIX: ICON - Gray "Auto Ammoloader" missing completely
FIX: ICON - Gray "MK-1 Machine Gun" missing completely
FIX: ICON - Some gray versions are slightly misaligned compared to the color version
FIX: ICON - Gray "Prototype SK Oculars" incorrectly mapped
FIX: ICON - Gray "Prototype SK Chassis" incorrectly mapped
FIX: ICON - Gray "Contruction Frame" incorrectly mapped
FIX: ICON - Gray "Any Distilled Spirits" not mapped
FIX: ICON - Gray "Leaf Fossil" not mapped
FIX: ICON - parcel_generic incorrectly mapped
FIX: EVENT - One of the fishpeople spawning events triggers an alert sound with no actual alert
FIX: EVENT - "Empire Times" script error
FIX: EVENT - "Creepy Night" causes many duplicates of the same memory
FIX: WORKSHOP - Mine shoring supplies consumed even on failed attempts
FIX: WORKSHOP - Premature "No Door" alert for the barracks
FIX: RECIPE - Caninha isn't flammable
FIX: RECIPE - "Wheat" missing all tags, impacts all recipes with it
FIX: RECIPE - research required - "Make Bottle of Sulphur Tonic" - chem/adv workbench - missing "Tonic Healing"
FIX: RECIPE - research required - "Make Bucket of Molasses" - kitchen/adv workbench - missing "Temperate/Desert"
FIX: RECIPE - "#x item needed" warning - "Roast Coffee Beans"
FIX: RECIPE - "#x item needed" warning - "Make Charcoal" - metal/brick char kiln
FIX: RECIPE - "#x item needed" warning - "Make MK1 Machine Gun" - metal/smithing forge
FIX: RECIPE - "#x item needed" warning - "Make Landmine" - metal/smith forge
FIX: RECIPE - "#x item needed" warning - "Brew Beer" - kitchen/brewing vat
FIX: RECIPE - "#x item needed" warning - "Make Prototype SK Chassis"
FIX: RECIPE - "#x item needed" warning - "Make Steam Knight Forge"
FIX: RECIPE - "#x item needed" warning - "Make Power Core Dynamo"
FIX: RECIPE - "#x item needed" warning - "Make Ammo Autoloader"
FIX: RECIPE - "#x item needed" warning - "Make Grenade Launcher Locker"
FIX: RECIPE - "#x item needed" warning - "Make Charcoal" - iron charcoal kiln
FIX: RECIPE - "#x item needed" warning - "Make Bucket of Lacquer" - chem workbench
FIX: RECIPE - "#x item needed" warning - "Make Reactive Catalyst" - chem workbench
FIX: RECIPE - "#x item needed" warning - "Make Reactive Catalyst" - adv workbench
FIX: RECIPE - "#x item needed" warning - "Make Prototype SK Oculars" - adv workbench
FIX: RECIPE - "#x item needed" warning - "Make Brick Charcoal Kiln" - adv workbench
FIX: RECIPE - "#x item needed" warning - "Make Dewatering Pump" - adv workbench
FIX: RECIPE - "#x item needed" warning - "Make Standing Desk" - adv workbench
FIX: RECIPE - "#x item needed" warning - "Make Bookshelf" - adv workbench
FIX: RECIPE - "#x item needed" warning - "Make Refined Food" - kitchen/iron oven
FIX: RECIPE - tooltip items/amounts - "Make Ammo Autoloader"
FIX: RECIPE - tooltip items/amounts - "Make SK-MK1 Chassis"
FIX: RECIPE - tooltip items/amounts - "Make Reactive Catalyst"
FIX: RECIPE - tooltip items/amounts - "Make Bucket of Lacquer"
FIX: RECIPE - tooltip items/amounts - "Make SK-MK1 Oculars"
FIX: RECIPE - tooltip items/amounts - "Make Jezail Rifle Locker"
FIX: RECIPE - tooltip items/amounts - "Make Carbine Locker"
FIX: RECIPE - tooltip items/amounts - "Make SK-MK1 Chassis"
FIX: RECIPE - tooltip items/amounts - "Make Grenade Launcher Locker"
FIX: RECIPE - tooltip items/amounts - "Make Power Core Dynamo"
FIX: RECIPE - tooltip items/amounts - "Make Tin Exotic Caviar"
FIX: RECIPE - tooltip items/amounts - "Make Prototype SK Oculars" - adv workbench
FIX: RECIPE - tooltip items/amounts - "Make Standing Desk" - adv workbench
FIX: RECIPE - tooltip items/amounts - "Make Bookshelf" - adv workbench
FIX: RECIPE - tooltip items/amounts - "Make Fancy Bookshelf" - adv workbench
FIX: RECIPE - tooltip items/amounts - "Make Charcoal" - metalworks/ind kiln
FIX: RECIPE - tooltip items/amounts - "Make Refined Food" - kitchen/iron oven
FIX: RECIPE - tooltip items/amounts - "Make SK-MK1 Oculars" - ceram/ceram workbench
FIX: RECIPE - actual items/amounts usage - "Make Ornate Bed" - too few Gold Ingots used
FIX: RECIPE - actual items/amounts usage - "Make Ornate Table and Chair Set" - too few Glass Panes used
FIX: RECIPE - actual items/amounts usage - "Make Brick Ceramics Kiln" - too few Bricks used
FIX: RECIPE - actual items/amounts usage - "Make Ceramics Press" - too many Brass Cogs used
FIX: RECIPE - actual items/amounts usage - "Make Stone Altar" - too few Stone used
FIX: RECIPE - actual items/amounts usage - "Make Iron Ceramics Kiln" - too few Iron Pipes used
FIX: RECIPE - actual items/amounts usage - "Make Jezail Rifle Locker" - Iron Ingots instead of Iron Plates used
FIX: RECIPE - actual items/amounts usage - "Make Stained Glass Window" - too few Bric-a-brac used
FIX: RECIPE - actual items/amounts usage - "Make Bucket of Lacquer" - adv workbench - too few yield
FIX: RECIPE - actual items/amounts usage - "Make Leyden Jars" - too many Copper Ingots used
FIX: RECIPE - actual items/amounts usage - "Make Cabbage Stew" - steam oven - wildly wrong cost
FIX: RECIPE - actual items/amounts usage - "Grind Stone Into Ore" - too many Stone used
FIX: TEXT - "Bricks (5)" references instead of "Bricks"
FIX: TEXT - "Steam Knight MK1" references instead of "MK-1 Steam Knight"
FIX: TEXT - Several "peice" vs "piece" typos
FIX: TEXT - Phased out "Basic Food" referenced instead of "Cooked Meat"
CHANGE: UI - More tooltips scale with the content, preventing huge (mostly empty) tooltips
CHANGE: ICON - "Any Timber" icon (color+gray) created
CHANGE: ICON - "Any Fuel" icon (color+gray) created
CHANGE: ICON - "Ore" category icon (color+gray) created
CHANGE: ICON - Recipes using "Any Distilled Spirits" now show the appropriate icon
CHANGE: ICON - Recipes using "Any Timber" now show the new icon
CHANGE: ICON - Recipes using "Any Fuel" now show the new icon
CHANGE: RECIPE - "Fungus" changed to "Any Fungus"
CHANGE: TEXT - Recipes referencing logs now state "Any Timber"
CHANGE: TEXT - Recipes referencing coal or charcoal now state "Any Fuel"
CHANGE: TEXT - Recipes referencing raw meat now state "Any Raw Meat"
CHANGE: TEXT - Recipes referencing vegetables now state "Any Vegetable"
CHANGE: TEXT - Recipes referencing fruit now state "Any Fruit"
CHANGE: EVENT - Some fishpeople spawning events have been given the chance to remain mysterious
CLEANUP: Removed redundant "Make Steam Oven" job code
CLEANUP: Removed redundant "Make Maize Chowder" job code




                                      ,
                                      ;\   It's my incessent
                                     /  \  droning, isn't it?
                   Please, God,      `.  ]          ,^^--.      Are we dead, Mike?
                   say "The End". __   [  \        /      \              ,',^-_
                                 /   \ !   \       |       |            / /   /
                                 \   /  \   \      |       ;      ,__. |    ,'
                                  < |    \   `.    |      /      (    `   __>
                                ,_| |_.   \    `-__>      >.      `---'\ /
                               /,.   ..\   `.               `.         | |
                               U |   | U     `.               \    ,--~   ~--.
--~~~~--_       _--~~~~--_       _--~~~~--_    \  _--~~~~--_   \  /_--~~~~~--_\
         `.   ,'          `.   ,'          `.  |,'          `.  \,'           `.
           \ /              \ /              \ /              \ /               \