# CMSIS Pack Index Generator Tool

CMSIS-Pack intro goes here

# How to use it

1. Add a vendor index file (link-to-vdix-page)

    $ cmpack add-vidx TheVendor https://the-vidx.com/TheVendor.vidx

2. Update the list of packages

    $ cmpack update

3. Search for packages or components or devices

    $ cmpack search pack-you-are-looking-for

4. Install a pack

    $ cmpack install TheVendor.ThePack[:.0.0.1]

5. Remove a pack

    $ cmpack remove TheVendor.ThePack[:.0.0.1]
