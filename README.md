#Overview

lev mintscan api to find validators who havent voted on a proposal and dump their twitter handles into blocks so they can be notified via twitter that they aren't voting

to execute, run:

```
go run main.go PROPOSALID
```

ex:

```
go run main.go 32
```

it will dump something like this. easily digestible lines of <= 280 character length:

```


@CosmostationVD @OmniFlixNetwork @swiss_staking @Sentinel_co @frensvalidator @stakecito @CephalopodEquip @MadeinBlock @kingnodes @ChorusOne @nodes_smart @KlubStaking @kalpa_tech @ChandraStation @ben2x4 @EZStaking @deuslabs @EarthValidator @BlueStakeNet @_elsehow @strangelovefund




@crypto_crew @rhinostake @blockpane @Figmentio @DHKdao @StakeLab @stakewolle @SecureSecrets @ValidatingChaos @HashQuark @ShapeShift @nacion_crypto @basblock_io @StakeEco @web34ever @CosmicValidator @01node @ValidatorRun @ChillinValidtn @cosmology_tech @chainflowpos




@ArtifactStaking @StakeTartare @chainlayerio @Catboss_Network @Cosmos_Spaces @SimplyStaking @BlockhuntersOrg @BlockdaemonHQ @AlterDapp @uGaenn @Commercionet @_D_Whale @Redline_Val


```
