
0.7.30
=============
2021-07-21

* Fix abug in the kubernetes create cmd, modify the name of the utils for check the length of the name (2d4641bf)

0.7.29
=============
2021-07-21

* Fixed error in the creation of the cluster (dc4d01b3)

0.7.28
=============
2021-07-21

* Added node as required flag in the scaling cmd (681149f7)
* Fixed error in the cmd add app to the cluster (8a37fada)

0.7.27
=============
2021-07-18

* Merge pull request #105 from carehart/patch-1 (8e04f2a7)
* Added more info to all error (5961cbf2)
* Updated the Changelog (c638815f)
* Update the civogo lib to the v0.2.48 (9c740a65)
* Updated changelog (68719517)
* Merge pull request #103 from civo/Fix-apikey-newline (3489c95a)
* Merge pull request #102 from civo/apikey-stdin (7da655a3)
* Merge branch 'master' of https://github.com/civo/cli (9c4ebe11)
* Updated Changelog (e0370def)
* Fixed error in the network rm command (117ebe6c)
* Updated the Changelog (144c031f)
* Fixed error at the moment to set a new region (1d924e09)
* updated the README.md (3c4f01a4)
* Updated the Changelog (91217a27)
* Standardized json outputs (bc005a62)
* Updated the Changelog (d6bdf67c)
* Updated the civogo lib to 0.2.47 (e768441c)
* Merge branch 'master' of https://github.com/civo/cli (4f6e53f8)
* Added changes in the json output for all cmd (32a6e079)
* Added the --pretty global flag to print the json in pretty format (52bdd806)
* Merge branch 'master' of https://github.com/civo/cli (a77bb430)
* Fix error in the application show cmd (c54fdf15)
* Remove debug statement (912d2ba9)
* Fix for breaking new customers adding their first key (0172994e)
* Updated the changelog (f58d3a02)
* Fixed error introduce fixing the issue #90 (1994ff59)
* Updated changelog (fd42c9f6)
* Fixed error in every command exit status zero (2d838d49)
* Merge pull request #89 from rajatjindal/fix-log-msg (dc724656)
* Updated the changelog (197dc49e)
* Updated the civogo lib to 0.2.45 how fix problem in the instance cmd (faff826d)
* Updated the Changelog (4acd29bc)
* Merge pull request #87 from civo/fix-default-region (273cdb01)
* Updated Changelog (64b63375)
* Fixed error in instance create cmd, now you can delete multiple pool in a cluster (8965abe8)
* Updated Changelog (670a4d02)
* Fixed kubernetes cmd and hidde volume for now (717623a6)
* Updated the changelog (7dc15bba)
* Added muti remove option to k3s, network, instance, domain, ssh key, firewall (2f2d4a9f)
* Fixed error in the show cmd in the instance (b779d6c5)
* Update README.md (d2737d89)
* Updated the changelog (bb4b1b50)
* Fixed error in the CLI also added new cmd to instance and kuberntes and new param to the size and updated the README (8d95e547)
* Added the las version of civogo lib v0.2.37 (a762c1d7)
* Updated the changelog (cec81869)
* Updated the civogo lib to 0.2.36 this will be fix scaling , rename bugs (52b2612a)
* Updated the changelog (8f7a1e40)
* Updated the civogo lib to v0.2.35 (5bde617f)
* Now domain create show nameserver (b6b7dd60)
* Fixed error in CheckAPPName func (4fe41126)
* Updated the changelog (4f930c71)
* fixed error in the output of the CLI (1e962709)
* Added verification to the kubernetes utils (a7e16d8e)
* Fixed more error in golang code (ed88a8a9)
* Added some verification to k8s cluster (3af72def)
* Updated the civogo lib to v0.2.32 (efdfff99)
* Merge pull request #72 from stiantoften/typofix (0b4e3570)
* Merge pull request #69 from DoNnMyTh/patch-1 (b686accc)
* Updated the changelog (49fece65)
* Added new cmd to show the post-install for every app installed in the cluster (48dbcaf2)
* Merge pull request #66 from DoNnMyTh/patch-4 (f851ed70)
* Merge pull request #62 from beret/darwin-arm64 (ef8e37ed)
* Updated the changelog (e0640948)
* Merge branch 'master' of https://github.com/civo/cli (ea7b5b89)
* Chnage to go 1.16 to build the binary for Apple M1 too (3156e008)
* Update the the civogo lib to 0.2.28 and fix error in the show command of kubernetes (ee1783fa)
* Merge pull request #55 from martynbristow/master (621e66a6)
* Updated the changelog (3f876df7)
* Fixed bug when you create a kubernetes cluster, also add suggest to the remove command (8346ff20)
* Fixed bug in the kubernetes create command (2ea5e0d6)
* Updated the changelog (b4d974b9)
* - Added the merge flag to the create command (9c259c32)
* Merge branch 'master' of https://github.com/civo/cli (0aa31637)
* Fix error in the creation of kubernertes (ad87c6f8)
* Merge branch 'master' of https://github.com/civo/cli (6ac46a0e)
* modified the network field to use the label value and not the name (69d94728)
* Updated the changelog (32dfda47)
* Added the `-network` param to the firewall cmd (8bcd181d)
* Updated the changelog (92905a29)
* Fixed error in the goreleaser file (57186e0c)
* Updated gorelease to remove some deprecated options (a8973554)
* Updated the civogo lib to v0.2.26 (63385a00)
* Updated the README.md (30439ca9)
* Merge branch 'master' of https://github.com/civo/cli (331e0289)
* Updated the Change log (915f932e)
* Added the option to set the region in the config file (b9179edb)
* Support CIVOCONFIG as an ENV variable to override config file location (aeab4452)
* Fixed error in the README (82e09be5)
* Updated the Change log (c0992d34)
* Fixed error in the kubernetes list command (a1c3a587)
* Typo and readability fixes (#40) (17007e30)
* Merge pull request #38 from kaihoffman/master (98f6850d)
* Updated the Change log (34d531a1)
* Fixed error handling in all commands (5e789fe3)
* Updated the Change log (2d0cac5b)
* Update the civogo lib to v0.2.23 (2b8f4d92)
* Merge branch 'master' of https://github.com/civo/cli (9355e877)
* Fixed bug in the install script (585b9d16)
* Updated the Change log (9931669e)
* Merge pull request #34 from martynbristow/master (af8c5e7d)
* Updated the Change log (a2f4a29a)
* Updated the civogo lib to v0.2.22 (d5695bff)
* Merge pull request #33 from kaihoffman/master (a7b6ef7e)
* Updated the Change log (7ac95c35)
* Added the recycle cmd to kubernetes (d2457212)
* Updated the Change log (aebfdca4)
* Fixed the permission in the civo conf file (f89b57e5)
* Updated the Change log and README.md (aa214676)
* Added powershel and fish to the completion cmd (3ce2c969)
* Updated the civogo lib to v0.2.21 (8820f312)
* Update the Change log (cd403dc6)
* Merge pull request #30 from Johannestegner/master (b59581ea)
* Updated the apikey output (f2d000b9)
* Update README.md (318d5988)
* Updated the Change log (a4abc2b0)
* Add arm64 to build list (52beb72e)
* Updated the Change log (a1bbc43f)
* Updated the civogo lib to v0.2.19 (b99f4258)
* Update the Change log (9847d8c7)
* Updated the civogo lib to v0.2.18 (9879243a)
* Update README.md (102dd5d6)
* Update the gitignore (b4cce868)
* Update the README.md (999d0895)
* Updated Change log (8a7b09a7)
* Allowed SRV record type in the DNS command (594884e1)
* Updated civogo lib to v0.2.17 (10e6d50d)
* Update the Change log (1d9b34a1)
* Added STOPPING status (422d3933)
* Updated Change log (22d81d1e)
* improved the custom output (98939932)
* Updated Change log (9e48a633)
* Fixed error in the custom output (cc05acff)
* Updated the README (c8815d73)
* Update the Change log (3900c15e)
* Fixed typo error in the Dockerfile (9d479dae)
* Updated Change log (30376d69)
* Fixed error in the Dockerfile (73e5e0d8)
* Update the Change log (6054efff)
* Updated the Dockerfile to add curl (e2d0d5f8)
* Update the Change log (9ed54ddd)
* Added kubectl inside the docker image (65c2c23c)
* Updated Change log (f7fbd5be)
* Quota command improved (381ab929)
* Added new show command to apikey (cb8c73d5)
* Updated Change log (3fde5e29)
* Update the civogo lib to v0.2.16 (8257fb46)
* Fixed the config generator (5a2efa5e)
* Remove windows 386 build from goreleaser (ab8d0bba)
* Added verification to all commands (80a158ea)
* Add a little polish for DEVELOPER.md (2c616de4)
* Added the develper file (3fefaffb)
* Update the Change log (fab73578)
* Name changed to label in the message of remove (add3b7f4)
* Updated the message in the remove command (c7c57d9d)
* Update Chnage log (fbbdff45)
* Update civogo lib to v0.2.15 (09b3bbe2)
* Added verification before delete any object (a9c4334f)
* Checking correct flag for Kubernetes readiness (e0f543c4)
* Update Change log (e3c8e654)
* Merge pull request #15 from ssmiller25/k8s-doc-upd (9e0767c7)
* Update the Change log (d94a61c5)
* Added new features when creating a kubernetes cluster (fb51ae8c)
* Update Change log (90daeee9)
* Fixed bug in output writing component (53ef3806)
* Update the change log (4ae0bf8d)
* Update the civogo to 0.2.14, fixed error in the intance create cmd (78e107de)
* Update the change log (a6a6b4df)
* Fixed error in the instance module and in the color utility (8738797c)
* Update goreleaser conf (87542cb8)
* Fix error for goreleaser (50c1433c)
* Merge pull request #12 from civo/feature/auto_update (02a956f8)
* Update Change log (c10a147b)
* Fixed color in the error message (a30bb996)
* Update Change log (5257cc27)
* Added verification at the moment of delete a snapshot (911cb634)
* Change the lib to add color to the CLI (d35889a5)
* Update Change log (70021fc8)
* Fixed error in the configuration of the kubeconfig (87326d25)
* Update Change log (e702aa80)
* Merge all color utility in one place (6dd3b564)
* Fixed error in the --merge option for the kubernetes config in windows (1094c066)
* Update the README.md (cfb1259e)
* Update the change log (f811d3e2)
* Update the change log (91a21bdd)
* Added the --save and --switch option to kubernetes create (950a432d)
* Fixed typo in the kubernetes config error message (bf097440)
* Added verification to kubernetes config cmd (7e4bb504)
* Fix error in kubernetes utils (f1f46612)
* Now if you use --switch with --merge, automatically the cli will change the context in kubernetes (8b9fad39)
* Update the change log (5307d998)
* Removed option egress as direction when creating a firewall rule (7baf88f4)
* Changed the words used to define the direction in the firewall rule (3df9c5b6)
* Added .editorconfig (cfb6ca22)
* Update the Change log (442e1776)
* Added verification to the firewall rule creation (a7fc88f9)
* Added new verification step to the kubernetes config (f77cf09e)
* Update the change log (5dc1bcac)
* Added the option to install multiple application at the same time in the Kubernetes Cluster (611120b1)
* Update the change log (d6cc1d42)
* Added CPU, RAM and SSD fields to Instance and Kubernetes CMD (6c5cfec7)
* Update the change log (e88016ef)
* Added new feature (56719687)
* Fixed error in the Makefile (feb35382)
* Merge pull request #8 from ssmiller25/selective-make (6caba09c)
* Chnage the name of the fordel in Homebrew Formula (e96b84a4)
* Merge pull request #7 from ssmiller25/k3s-docs (4a70e48b)
* Fixed error in cli (10196c99)
* Add KubernetesVesion and UpgradeAvailableTo (d61e6a9a)
* Add nfpm and use the v2 of the goreleaser action (49d55069)
* Fixed the message to show after upgrade happened (396934a5)
* Upgrade the version of civogo lib (6ff45d0d)
* Fixed error in the kubernetes cmd (723028c1)
* Change the dockerfile to run as root (7f15c77b)
* fix: Fix error adding initil user (28c99ac9)
* Update the CHANGELOG (6669b296)
* fix: Fix message in the kubernetes upgrade cmd (05960c0b)
* docs: Update the CHANGELOG (ea04cc57)
* fix: Remove update from the cmd kubernetes show (6cb4144a)
* Added the CHANGELOG.md (e99843a4)
* Added new feature (f1726759)
* Fix correct version showing (83fe22c8)
* Added some feature (a356acaa)
* Fixed some errors (911f046e)
* Remove snap from release process (15d43a7c)
* Update kubernetes config (6413a22a)
* Update .goreleaser.yml (7ff511dd)
* Update .goreleaser.yml (a763c30b)
* Update .goreleaser.yml (3166ee2f)
* Fix error (787ad7de)
* Added some corrections (e5fe1faa)
* Fixed error in goreleaser (72708d8c)
* Modifications in logic (497e0750)
* Tweaks to install.sh (29250021)
* Add configuration to snap (148a0473)
* Update README.md (3bea636c)
* Update the badge (d118e1e3)
* Add badge to the README (0da2f4ff)
* Fix error in the config (8a1430f9)
* Fix user in the Dockerfile (b0d63fe2)
* Update the image to alpine the Dockerfile (0b182a0b)
* Update the Dockerfile (00942fdb)
* Update the cli (4d4ac165)
* Fix wording (0104cacd)
* Update README now this is the main tool (1597d82f)
* Remove comments (f42b14ac)
* Run new test (578d62d7)
* Test brews releaser (425f84f5)
* Try to fix Homebrew releasing (bdc4a966)
* Fix goreleaser.yml (4b0a350b)
* Add homebrew tap updating to goreleaser (2801d339)
* fix(Project): Added some fix (8ed9c0bf)
* Lots of minor tweaks (ef039b50)
* Changes to instance create (82901093)
* feat(Utility): Add TrackTime (7e16e37d)
* feat(Utility): Add random (268021cf)
* fix(Bug): Fix bug found in all code (47c14806)
* fix(Domain): Fix some error (04c0775e)
* feat(Install): Added install bash (2e584464)
* feat(Completion): Added completion cmd (11f44f4b)
* fix(Project): Fix the some bug (a082fdaf)
* fix(Volume): Fix the example (d3e46926)
* refactor(Project): Added example to all cmd (0119e351)
* fix(Network): Remove Name from the table (b5659302)
* fix(LoadBalancer): added new alias (2cb08795)
* refactor(Project): Modify all project error handler (36f932ee)
* fix(Version): fix text in version cmd (d4b80553)
* Merge remote-tracking branch 'origin/master' (732fcc41)
* fix(Version): fix text in version cmd (fe625b9b)
* fix(Version): update the version cmd (b4af168b)
* fix(Version): change the version cmd (17402302)
* fix(Version): update the version cmd (75fe68e8)
* fix(Version): Fixed the version cmd (f01cff84)
* fix(Goreleser): Fix the conf file (3ecf80b0)
* fix(Version): Fixed the version cmd (b35e0cf5)
* fix(Goreleser): Fix the conf file (122998f8)
* fix(Version): Upgrade the version of the civogo lib (c84e2461)
* feat(Version): Added the version cmd (aaa6a091)
* fix(Goreleser): Fix the conf file (f4b78e1b)
* fix(Goreleser): Fix the conf file (6beb0caf)
* fix(Goreleser): Fix the conf file (0953e531)
* fix(Goreleser): Fix the conf file (0ab428ae)
* fix(Goreleser): Fix the action file (2d05ff91)
* fix(Goreleser): Fix the conf file (5b12fdf7)
* fix(Goreleser): Fix the action file (44ad7dff)
* refactor(Project): Modify all project (bc2f508f)
* Merge pull request #2 from civo/dev (daccd7ba)
* Add all but one instance command (7ad58ed4)
* Split more commands out in to their own files (3f3b06e3)
* Split instance commands out to invidual files (2ce29790)
* Add console and reboot commands (d24951af)
* Add github token path to goreleaser (cd811766)
* Add testing workflow (d01e3ed8)
* Add goreleaser support (6410d22b)
* Add the first action command to instances (ec0634b6)
* Change OutputWriter to use labels as well as keys (52eeee28)
* Change to use Find* in civogo (e15de86c)
* Refactor to new OutputWriter (417c47aa)
* Add size command (f7201a09)
* Add quota command (a07f44a6)
* Add Region command (48cd4f0c)
* Complete API key management (bc9e3d53)
* Remove unnecessary comment (7f7fb17f)
* Remove (old) accidentally commited API keys (ea215e4f)
* Complete API key management functionality (94421615)
* Add API Key management (4beb5f1e)
* Update links in README (17848b68)
* Add plan to README (b9fc702c)
* First commit - initial makefile (38300f4e)
* feat(Template): Added the kubernetes cmd (0370dde0)
* feat(Goreleser): Fix the goreleaser file (59a4231c)
* feat(Template): Added the template cmd (4f406582)
* feat(Volume): Added the volume cmd (6a84a030)
* feat(snapshot): Added the snapshot cmd (775885d0)
* feat(network): Added the network cmd (6c6d035e)
* feat(sshkey): Added the ssh key cmd (4f3a6445)
* feat(loadbalancer): Added the cmd loadbalancer (410abb51)
* feat(firewall): Added the cmd firewall, and firewall rule (0ad42aed)
* feat(domain): Added the cmd domain record (51de188d)
* feat(domain): Added the cmd domain (698df535)
* feat(instance): Added the option to create instances (22dcec05)
* refactor(apikey): Modify the apikey (1d113038)
* refactor(apikey): Modify the apikey (430a4114)
* refactor(apikey): Modify the apikey function (a8310907)
* Update app name in tap-release (de91a6e0)
* Create automatic tap release (fa017dff)
* k8s usage update, and markdown linting (69f2d1d8)
* Selective Makefile (e6ed4e78)
* Updated the Change log (46b80c81)
* Added the auto update option to the cli (eb9ff675)
* Update on apps, minor quote addition (dc14aa15)
* Removed windows specific file removal, as it's not needed. (03287d31)
* Updated kubeconfig -s -m command to make sure the file is closed before trying to remove it. (2420613f)
* Grammar fixes in user output (d6ac0b6b)
* Fix Transcribe in WriteCustomOutput (813c6193)
* Failing Example Test for Fn WriteCustomOutput (0cc50313)
* Merge pull request #35 from jaysonsantos/patch-1 (7d637058)
* Add windows setup guide (d61ce9c2)
* Add recycle node info to README (910390c5)
* Reword error message when setting a region that doesn't exist (d872c509)
* Merge pull request #45 from DoNnMyTh/master (3d415046)
* Fixed typos, merged 2 region and regions section into one, added command to change region (c0982c82)
* Merge pull request #49 from DoNnMyTh/patch-1 (b3eacf1e)
* updated README.md (cc1bd1a5)
* Resolve Issues with install.sh (bef90aa0)
* Resolve missing fatal function (4cc77f89)
* Merge pull request #59 from DoNnMyTh/patch-2 (b4c36996)
* Merge pull request #58 from shanmukhkotharu/master (30461e9f)
* Update README.md (ec43b0a6)
* Update README.md (0fd2e34a)
* Updated Readme.md for chocolatey support (ad84c28d)
* Use arm64 binaries for arm64/aarch64 platforms (827c7d65)
* Minor README formatting for macOS and Civo CLI (2dadf24b)
* Updated readme.md (e0043d4a)
* Fixed spelling mistake (d17b03f9)
* Fix typo in default config (b1e0fc94)
* Add missing default region warning to all commands that need one (48f30857)
* Remove global warning on default regions as some commands don't need them (f91e4211)
* Add utility method to check for if a current region is set (36399b76)
* Don't set a default region when creating blank configs (and certainly not SVG1) (8986e650)
* Fix typo in variable name (824b617e)
* Get the default region from the API when adding an API key, if one isn't set (f9df4cd2)
* fix log msg (e3ca6605)
* Update README.md (f513044c)
* Merge pull request #95 from SubhasmitaSw/documentation (8050fae9)
* Update DEVELOPER.md (8c9aeadb)
* requested changes are updated. (fe316872)
* Improved the documentation and fixed little typos and punctuations (455828a0)
* Move error check (f3ea046b)
* Merge pull request #99 from dirien/fix-84 (f9476051)
* Merge branch 'master' of https://github.com/civo/cli (d96aa357)
* Update changelog (1b03bb2a)
* Add support to save API key from stdin (182692ca)
* Update changelog (7a9c9be1)
* Fix: removal of newline for api key name (d3533c97)
* correct grammar mistake in console message (af2fbafd)

0.7.26
=============
2021-07-15

* Update the civogo lib to the v0.2.48 (9c740a65)
* Updated changelog (68719517)
* Merge pull request #103 from civo/Fix-apikey-newline (3489c95a)
* Merge pull request #102 from civo/apikey-stdin (7da655a3)
* Merge branch 'master' of https://github.com/civo/cli (9c4ebe11)
* Updated Changelog (e0370def)
* Fixed error in the network rm command (117ebe6c)
* Updated the Changelog (144c031f)
* Fixed error at the moment to set a new region (1d924e09)
* updated the README.md (3c4f01a4)
* Updated the Changelog (91217a27)
* Standardized json outputs (bc005a62)
* Updated the Changelog (d6bdf67c)
* Updated the civogo lib to 0.2.47 (e768441c)
* Merge branch 'master' of https://github.com/civo/cli (4f6e53f8)
* Added changes in the json output for all cmd (32a6e079)
* Added the --pretty global flag to print the json in pretty format (52bdd806)
* Merge branch 'master' of https://github.com/civo/cli (a77bb430)
* Fix error in the application show cmd (c54fdf15)
* Remove debug statement (912d2ba9)
* Fix for breaking new customers adding their first key (0172994e)
* Updated the changelog (f58d3a02)
* Fixed error introduce fixing the issue #90 (1994ff59)
* Updated changelog (fd42c9f6)
* Fixed error in every command exit status zero (2d838d49)
* Merge pull request #89 from rajatjindal/fix-log-msg (dc724656)
* Updated the changelog (197dc49e)
* Updated the civogo lib to 0.2.45 how fix problem in the instance cmd (faff826d)
* Updated the Changelog (4acd29bc)
* Merge pull request #87 from civo/fix-default-region (273cdb01)
* Updated Changelog (64b63375)
* Fixed error in instance create cmd, now you can delete multiple pool in a cluster (8965abe8)
* Updated Changelog (670a4d02)
* Fixed kubernetes cmd and hidde volume for now (717623a6)
* Updated the changelog (7dc15bba)
* Added muti remove option to k3s, network, instance, domain, ssh key, firewall (2f2d4a9f)
* Fixed error in the show cmd in the instance (b779d6c5)
* Update README.md (d2737d89)
* Updated the changelog (bb4b1b50)
* Fixed error in the CLI also added new cmd to instance and kuberntes and new param to the size and updated the README (8d95e547)
* Added the las version of civogo lib v0.2.37 (a762c1d7)
* Updated the changelog (cec81869)
* Updated the civogo lib to 0.2.36 this will be fix scaling , rename bugs (52b2612a)
* Updated the changelog (8f7a1e40)
* Updated the civogo lib to v0.2.35 (5bde617f)
* Now domain create show nameserver (b6b7dd60)
* Fixed error in CheckAPPName func (4fe41126)
* Updated the changelog (4f930c71)
* fixed error in the output of the CLI (1e962709)
* Added verification to the kubernetes utils (a7e16d8e)
* Fixed more error in golang code (ed88a8a9)
* Added some verification to k8s cluster (3af72def)
* Updated the civogo lib to v0.2.32 (efdfff99)
* Merge pull request #72 from stiantoften/typofix (0b4e3570)
* Merge pull request #69 from DoNnMyTh/patch-1 (b686accc)
* Updated the changelog (49fece65)
* Added new cmd to show the post-install for every app installed in the cluster (48dbcaf2)
* Merge pull request #66 from DoNnMyTh/patch-4 (f851ed70)
* Merge pull request #62 from beret/darwin-arm64 (ef8e37ed)
* Updated the changelog (e0640948)
* Merge branch 'master' of https://github.com/civo/cli (ea7b5b89)
* Chnage to go 1.16 to build the binary for Apple M1 too (3156e008)
* Update the the civogo lib to 0.2.28 and fix error in the show command of kubernetes (ee1783fa)
* Merge pull request #55 from martynbristow/master (621e66a6)
* Updated the changelog (3f876df7)
* Fixed bug when you create a kubernetes cluster, also add suggest to the remove command (8346ff20)
* Fixed bug in the kubernetes create command (2ea5e0d6)
* Updated the changelog (b4d974b9)
* - Added the merge flag to the create command (9c259c32)
* Merge branch 'master' of https://github.com/civo/cli (0aa31637)
* Fix error in the creation of kubernertes (ad87c6f8)
* Merge branch 'master' of https://github.com/civo/cli (6ac46a0e)
* modified the network field to use the label value and not the name (69d94728)
* Updated the changelog (32dfda47)
* Added the `-network` param to the firewall cmd (8bcd181d)
* Updated the changelog (92905a29)
* Fixed error in the goreleaser file (57186e0c)
* Updated gorelease to remove some deprecated options (a8973554)
* Updated the civogo lib to v0.2.26 (63385a00)
* Updated the README.md (30439ca9)
* Merge branch 'master' of https://github.com/civo/cli (331e0289)
* Updated the Change log (915f932e)
* Added the option to set the region in the config file (b9179edb)
* Support CIVOCONFIG as an ENV variable to override config file location (aeab4452)
* Fixed error in the README (82e09be5)
* Updated the Change log (c0992d34)
* Fixed error in the kubernetes list command (a1c3a587)
* Typo and readability fixes (#40) (17007e30)
* Merge pull request #38 from kaihoffman/master (98f6850d)
* Updated the Change log (34d531a1)
* Fixed error handling in all commands (5e789fe3)
* Updated the Change log (2d0cac5b)
* Update the civogo lib to v0.2.23 (2b8f4d92)
* Merge branch 'master' of https://github.com/civo/cli (9355e877)
* Fixed bug in the install script (585b9d16)
* Updated the Change log (9931669e)
* Merge pull request #34 from martynbristow/master (af8c5e7d)
* Updated the Change log (a2f4a29a)
* Updated the civogo lib to v0.2.22 (d5695bff)
* Merge pull request #33 from kaihoffman/master (a7b6ef7e)
* Updated the Change log (7ac95c35)
* Added the recycle cmd to kubernetes (d2457212)
* Updated the Change log (aebfdca4)
* Fixed the permission in the civo conf file (f89b57e5)
* Updated the Change log and README.md (aa214676)
* Added powershel and fish to the completion cmd (3ce2c969)
* Updated the civogo lib to v0.2.21 (8820f312)
* Update the Change log (cd403dc6)
* Merge pull request #30 from Johannestegner/master (b59581ea)
* Updated the apikey output (f2d000b9)
* Update README.md (318d5988)
* Updated the Change log (a4abc2b0)
* Add arm64 to build list (52beb72e)
* Updated the Change log (a1bbc43f)
* Updated the civogo lib to v0.2.19 (b99f4258)
* Update the Change log (9847d8c7)
* Updated the civogo lib to v0.2.18 (9879243a)
* Update README.md (102dd5d6)
* Update the gitignore (b4cce868)
* Update the README.md (999d0895)
* Updated Change log (8a7b09a7)
* Allowed SRV record type in the DNS command (594884e1)
* Updated civogo lib to v0.2.17 (10e6d50d)
* Update the Change log (1d9b34a1)
* Added STOPPING status (422d3933)
* Updated Change log (22d81d1e)
* improved the custom output (98939932)
* Updated Change log (9e48a633)
* Fixed error in the custom output (cc05acff)
* Updated the README (c8815d73)
* Update the Change log (3900c15e)
* Fixed typo error in the Dockerfile (9d479dae)
* Updated Change log (30376d69)
* Fixed error in the Dockerfile (73e5e0d8)
* Update the Change log (6054efff)
* Updated the Dockerfile to add curl (e2d0d5f8)
* Update the Change log (9ed54ddd)
* Added kubectl inside the docker image (65c2c23c)
* Updated Change log (f7fbd5be)
* Quota command improved (381ab929)
* Added new show command to apikey (cb8c73d5)
* Updated Change log (3fde5e29)
* Update the civogo lib to v0.2.16 (8257fb46)
* Fixed the config generator (5a2efa5e)
* Remove windows 386 build from goreleaser (ab8d0bba)
* Added verification to all commands (80a158ea)
* Add a little polish for DEVELOPER.md (2c616de4)
* Added the develper file (3fefaffb)
* Update the Change log (fab73578)
* Name changed to label in the message of remove (add3b7f4)
* Updated the message in the remove command (c7c57d9d)
* Update Chnage log (fbbdff45)
* Update civogo lib to v0.2.15 (09b3bbe2)
* Added verification before delete any object (a9c4334f)
* Checking correct flag for Kubernetes readiness (e0f543c4)
* Update Change log (e3c8e654)
* Merge pull request #15 from ssmiller25/k8s-doc-upd (9e0767c7)
* Update the Change log (d94a61c5)
* Added new features when creating a kubernetes cluster (fb51ae8c)
* Update Change log (90daeee9)
* Fixed bug in output writing component (53ef3806)
* Update the change log (4ae0bf8d)
* Update the civogo to 0.2.14, fixed error in the intance create cmd (78e107de)
* Update the change log (a6a6b4df)
* Fixed error in the instance module and in the color utility (8738797c)
* Update goreleaser conf (87542cb8)
* Fix error for goreleaser (50c1433c)
* Merge pull request #12 from civo/feature/auto_update (02a956f8)
* Update Change log (c10a147b)
* Fixed color in the error message (a30bb996)
* Update Change log (5257cc27)
* Added verification at the moment of delete a snapshot (911cb634)
* Change the lib to add color to the CLI (d35889a5)
* Update Change log (70021fc8)
* Fixed error in the configuration of the kubeconfig (87326d25)
* Update Change log (e702aa80)
* Merge all color utility in one place (6dd3b564)
* Fixed error in the --merge option for the kubernetes config in windows (1094c066)
* Update the README.md (cfb1259e)
* Update the change log (f811d3e2)
* Update the change log (91a21bdd)
* Added the --save and --switch option to kubernetes create (950a432d)
* Fixed typo in the kubernetes config error message (bf097440)
* Added verification to kubernetes config cmd (7e4bb504)
* Fix error in kubernetes utils (f1f46612)
* Now if you use --switch with --merge, automatically the cli will change the context in kubernetes (8b9fad39)
* Update the change log (5307d998)
* Removed option egress as direction when creating a firewall rule (7baf88f4)
* Changed the words used to define the direction in the firewall rule (3df9c5b6)
* Added .editorconfig (cfb6ca22)
* Update the Change log (442e1776)
* Added verification to the firewall rule creation (a7fc88f9)
* Added new verification step to the kubernetes config (f77cf09e)
* Update the change log (5dc1bcac)
* Added the option to install multiple application at the same time in the Kubernetes Cluster (611120b1)
* Update the change log (d6cc1d42)
* Added CPU, RAM and SSD fields to Instance and Kubernetes CMD (6c5cfec7)
* Update the change log (e88016ef)
* Added new feature (56719687)
* Fixed error in the Makefile (feb35382)
* Merge pull request #8 from ssmiller25/selective-make (6caba09c)
* Chnage the name of the fordel in Homebrew Formula (e96b84a4)
* Merge pull request #7 from ssmiller25/k3s-docs (4a70e48b)
* Fixed error in cli (10196c99)
* Add KubernetesVesion and UpgradeAvailableTo (d61e6a9a)
* Add nfpm and use the v2 of the goreleaser action (49d55069)
* Fixed the message to show after upgrade happened (396934a5)
* Upgrade the version of civogo lib (6ff45d0d)
* Fixed error in the kubernetes cmd (723028c1)
* Change the dockerfile to run as root (7f15c77b)
* fix: Fix error adding initil user (28c99ac9)
* Update the CHANGELOG (6669b296)
* fix: Fix message in the kubernetes upgrade cmd (05960c0b)
* docs: Update the CHANGELOG (ea04cc57)
* fix: Remove update from the cmd kubernetes show (6cb4144a)
* Added the CHANGELOG.md (e99843a4)
* Added new feature (f1726759)
* Fix correct version showing (83fe22c8)
* Added some feature (a356acaa)
* Fixed some errors (911f046e)
* Remove snap from release process (15d43a7c)
* Update kubernetes config (6413a22a)
* Update .goreleaser.yml (7ff511dd)
* Update .goreleaser.yml (a763c30b)
* Update .goreleaser.yml (3166ee2f)
* Fix error (787ad7de)
* Added some corrections (e5fe1faa)
* Fixed error in goreleaser (72708d8c)
* Modifications in logic (497e0750)
* Tweaks to install.sh (29250021)
* Add configuration to snap (148a0473)
* Update README.md (3bea636c)
* Update the badge (d118e1e3)
* Add badge to the README (0da2f4ff)
* Fix error in the config (8a1430f9)
* Fix user in the Dockerfile (b0d63fe2)
* Update the image to alpine the Dockerfile (0b182a0b)
* Update the Dockerfile (00942fdb)
* Update the cli (4d4ac165)
* Fix wording (0104cacd)
* Update README now this is the main tool (1597d82f)
* Remove comments (f42b14ac)
* Run new test (578d62d7)
* Test brews releaser (425f84f5)
* Try to fix Homebrew releasing (bdc4a966)
* Fix goreleaser.yml (4b0a350b)
* Add homebrew tap updating to goreleaser (2801d339)
* fix(Project): Added some fix (8ed9c0bf)
* Lots of minor tweaks (ef039b50)
* Changes to instance create (82901093)
* feat(Utility): Add TrackTime (7e16e37d)
* feat(Utility): Add random (268021cf)
* fix(Bug): Fix bug found in all code (47c14806)
* fix(Domain): Fix some error (04c0775e)
* feat(Install): Added install bash (2e584464)
* feat(Completion): Added completion cmd (11f44f4b)
* fix(Project): Fix the some bug (a082fdaf)
* fix(Volume): Fix the example (d3e46926)
* refactor(Project): Added example to all cmd (0119e351)
* fix(Network): Remove Name from the table (b5659302)
* fix(LoadBalancer): added new alias (2cb08795)
* refactor(Project): Modify all project error handler (36f932ee)
* fix(Version): fix text in version cmd (d4b80553)
* Merge remote-tracking branch 'origin/master' (732fcc41)
* fix(Version): fix text in version cmd (fe625b9b)
* fix(Version): update the version cmd (b4af168b)
* fix(Version): change the version cmd (17402302)
* fix(Version): update the version cmd (75fe68e8)
* fix(Version): Fixed the version cmd (f01cff84)
* fix(Goreleser): Fix the conf file (3ecf80b0)
* fix(Version): Fixed the version cmd (b35e0cf5)
* fix(Goreleser): Fix the conf file (122998f8)
* fix(Version): Upgrade the version of the civogo lib (c84e2461)
* feat(Version): Added the version cmd (aaa6a091)
* fix(Goreleser): Fix the conf file (f4b78e1b)
* fix(Goreleser): Fix the conf file (6beb0caf)
* fix(Goreleser): Fix the conf file (0953e531)
* fix(Goreleser): Fix the conf file (0ab428ae)
* fix(Goreleser): Fix the action file (2d05ff91)
* fix(Goreleser): Fix the conf file (5b12fdf7)
* fix(Goreleser): Fix the action file (44ad7dff)
* refactor(Project): Modify all project (bc2f508f)
* Merge pull request #2 from civo/dev (daccd7ba)
* Add all but one instance command (7ad58ed4)
* Split more commands out in to their own files (3f3b06e3)
* Split instance commands out to invidual files (2ce29790)
* Add console and reboot commands (d24951af)
* Add github token path to goreleaser (cd811766)
* Add testing workflow (d01e3ed8)
* Add goreleaser support (6410d22b)
* Add the first action command to instances (ec0634b6)
* Change OutputWriter to use labels as well as keys (52eeee28)
* Change to use Find* in civogo (e15de86c)
* Refactor to new OutputWriter (417c47aa)
* Add size command (f7201a09)
* Add quota command (a07f44a6)
* Add Region command (48cd4f0c)
* Complete API key management (bc9e3d53)
* Remove unnecessary comment (7f7fb17f)
* Remove (old) accidentally commited API keys (ea215e4f)
* Complete API key management functionality (94421615)
* Add API Key management (4beb5f1e)
* Update links in README (17848b68)
* Add plan to README (b9fc702c)
* First commit - initial makefile (38300f4e)
* feat(Template): Added the kubernetes cmd (0370dde0)
* feat(Goreleser): Fix the goreleaser file (59a4231c)
* feat(Template): Added the template cmd (4f406582)
* feat(Volume): Added the volume cmd (6a84a030)
* feat(snapshot): Added the snapshot cmd (775885d0)
* feat(network): Added the network cmd (6c6d035e)
* feat(sshkey): Added the ssh key cmd (4f3a6445)
* feat(loadbalancer): Added the cmd loadbalancer (410abb51)
* feat(firewall): Added the cmd firewall, and firewall rule (0ad42aed)
* feat(domain): Added the cmd domain record (51de188d)
* feat(domain): Added the cmd domain (698df535)
* feat(instance): Added the option to create instances (22dcec05)
* refactor(apikey): Modify the apikey (1d113038)
* refactor(apikey): Modify the apikey (430a4114)
* refactor(apikey): Modify the apikey function (a8310907)
* Update app name in tap-release (de91a6e0)
* Create automatic tap release (fa017dff)
* k8s usage update, and markdown linting (69f2d1d8)
* Selective Makefile (e6ed4e78)
* Updated the Change log (46b80c81)
* Added the auto update option to the cli (eb9ff675)
* Update on apps, minor quote addition (dc14aa15)
* Removed windows specific file removal, as it's not needed. (03287d31)
* Updated kubeconfig -s -m command to make sure the file is closed before trying to remove it. (2420613f)
* Grammar fixes in user output (d6ac0b6b)
* Fix Transcribe in WriteCustomOutput (813c6193)
* Failing Example Test for Fn WriteCustomOutput (0cc50313)
* Merge pull request #35 from jaysonsantos/patch-1 (7d637058)
* Add windows setup guide (d61ce9c2)
* Add recycle node info to README (910390c5)
* Reword error message when setting a region that doesn't exist (d872c509)
* Merge pull request #45 from DoNnMyTh/master (3d415046)
* Fixed typos, merged 2 region and regions section into one, added command to change region (c0982c82)
* Merge pull request #49 from DoNnMyTh/patch-1 (b3eacf1e)
* updated README.md (cc1bd1a5)
* Resolve Issues with install.sh (bef90aa0)
* Resolve missing fatal function (4cc77f89)
* Merge pull request #59 from DoNnMyTh/patch-2 (b4c36996)
* Merge pull request #58 from shanmukhkotharu/master (30461e9f)
* Update README.md (ec43b0a6)
* Update README.md (0fd2e34a)
* Updated Readme.md for chocolatey support (ad84c28d)
* Use arm64 binaries for arm64/aarch64 platforms (827c7d65)
* Minor README formatting for macOS and Civo CLI (2dadf24b)
* Updated readme.md (e0043d4a)
* Fixed spelling mistake (d17b03f9)
* Fix typo in default config (b1e0fc94)
* Add missing default region warning to all commands that need one (48f30857)
* Remove global warning on default regions as some commands don't need them (f91e4211)
* Add utility method to check for if a current region is set (36399b76)
* Don't set a default region when creating blank configs (and certainly not SVG1) (8986e650)
* Fix typo in variable name (824b617e)
* Get the default region from the API when adding an API key, if one isn't set (f9df4cd2)
* fix log msg (e3ca6605)
* Update README.md (f513044c)
* Merge pull request #95 from SubhasmitaSw/documentation (8050fae9)
* Update DEVELOPER.md (8c9aeadb)
* requested changes are updated. (fe316872)
* Improved the documentation and fixed little typos and punctuations (455828a0)
* Move error check (f3ea046b)
* Merge pull request #99 from dirien/fix-84 (f9476051)
* Merge branch 'master' of https://github.com/civo/cli (d96aa357)
* Update changelog (1b03bb2a)
* Add support to save API key from stdin (182692ca)
* Update changelog (7a9c9be1)
* Fix: removal of newline for api key name (d3533c97)

0.7.25
=============
2021-07-09

* Merge pull request #103 from civo/Fix-apikey-newline (3489c95a)
* Merge pull request #102 from civo/apikey-stdin (7da655a3)

0.7.24
=============
2021-07-01

* Merge branch 'master' of https://github.com/civo/cli (d96aa357)
* Fixed error in the network rm command (117ebe6c)

0.7.23
=============
2021-06-11

* Fixed error at the moment to set a new region (1d924e09)
* updated the README.md (3c4f01a4)

0.7.22
=============
2021-05-23

* Standardized json outputs (bc005a62)

0.7.21
=============
2021-05-23

* Updated the civogo lib to 0.2.47 (e768441c)
* Merge branch 'master' of https://github.com/civo/cli (4f6e53f8)
* Added changes in the json output for all cmd (32a6e079)
* Added the --pretty global flag to print the json in pretty format (52bdd806)
* Merge branch 'master' of https://github.com/civo/cli (a77bb430)
* Fix error in the application show cmd (c54fdf15)
* Remove debug statement (912d2ba9)
* Fix for breaking new customers adding their first key (0172994e)

0.7.19
=============
2021-05-11

* Fixed error introduce fixing the issue #90 (1994ff59)

0.7.18
=============
2021-05-11

* Fixed error in every command exit status zero (2d838d49)
* Merge pull request #89 from rajatjindal/fix-log-msg (dc724656)

0.7.17
=============
2021-05-10

* Updated the civogo lib to 0.2.45 how fix problem in the instance cmd (faff826d)

0.7.16
=============
2021-05-10

* Add missing default region warning to all commands that need one (48f30857)
* Remove global warning on default regions as some commands don't need them (f91e4211)
* Add utility method to check for if a current region is set (36399b76)
* Don't set a default region when creating blank configs (and certainly not SVG1) (8986e650)
* Fix typo in variable name (824b617e)
* Get the default region from the API when adding an API key, if one isn't set (f9df4cd2)

0.7.15
=============
2021-05-08

* Fixed error in instance create cmd, now you can delete multiple pool in a cluster (8965abe8)

0.7.14
=============
2021-05-06

* Fixed kubernetes cmd and hidde volume for now (717623a6)

0.7.13
=============
2021-04-19

* Added multi remove option to k3s, network, instance, domain, ssh key, firewall (2f2d4a9f)
* Fixed error in the show cmd in the instance (b779d6c5)
* Update README.md (d2737d89)

0.7.12
=============
2021-04-14

* Fixed error in the CLI also added new cmd to instance and kuberntes and new param to the size and updated the README (8d95e547)
* Added the last version of civogo lib v0.2.37 (a762c1d7)

0.7.11
=============
2021-04-06

* Updated the civogo lib to 0.2.36 this will be fix scaling , rename bugs (52b2612a)

0.7.10
=============
2021-04-02

* Updated the civogo lib to v0.2.35 (5bde617f)
* Now domain create show nameserver (b6b7dd60)
* Fixed error in CheckAPPName func (4fe41126)

0.7.9
=============
2021-03-12

* Fixed error in the output of the CLI (1e962709)
* Added verification to the kubernetes utils (a7e16d8e)
* Fixed more error in golang code (ed88a8a9)
* Added some verification to k8s cluster (3af72def)
* Updated the civogo lib to v0.2.32 (efdfff99)
* Merge pull request #72 from stiantoften/typofix (0b4e3570)
* Merge pull request #69 from DoNnMyTh/patch-1 (b686accc)

0.7.7
=============
2021-03-10

* Added new cmd to show the post-install for every app installed in the cluster (48dbcaf2)
* Merge pull request #66 from DoNnMyTh/patch-4 (f851ed70)
* Merge pull request #62 from beret/darwin-arm64 (ef8e37ed)

0.7.6
=============
2021-02-28

* Merge branch 'master' of https://github.com/civo/cli (ea7b5b89)
* Chnage to go 1.16 to build the binary for Apple M1 too (3156e008)
* Update the the civogo lib to 0.2.28 and fix error in the show command of kubernetes (ee1783fa)
* Merge pull request #55 from martynbristow/master (621e66a6)
* Updated the changelog (3f876df7)

0.7.4
=============
2021-02-20

* Fixed bug when you create a kubernetes cluster, also add suggest to the remove command (8346ff20)
* Fixed bug in the kubernetes create command (2ea5e0d6)
* Updated the changelog (b4d974b9)

0.7.3
=============
2021-02-19

* - Added the merge flag to the create command (9c259c32)
* Merge branch 'master' of https://github.com/civo/cli (0aa31637)
* Fix error in the creation of kubernertes (ad87c6f8)
* Merge branch 'master' of https://github.com/civo/cli (6ac46a0e)
* modified the network field to use the label value and not the name (69d94728)
* Updated the changelog (32dfda47)

0.7.2
=============
2021-02-10

* Added the `-network` param to the firewall cmd also added the network to the list of firewall and updated to the last version of the civogo lib v0.2.27 (8bcd181d)
* Updated the changelog (92905a29)

0.7.1
=============
2021-02-09

* Fixed error in the goreleaser file (57186e0c)
* Updated gorelease to remove some deprecated options (a8973554)
* Updated the civogo lib to v0.2.26, Added a check in the instance and kubernetes command and remove some field from the region list (63385a00)
* Updated the README.md (30439ca9)

0.7.0
=============
2021-02-08

* Added region to the CLI (cc0a3c2a)
* Added the option to set the region in the config file (b9179edb)
* Support CIVOCONFIG as an ENV variable to override config file location (aeab4452)
* Fixed error in the README (82e09be5)

0.6.46
=============
2020-12-10

* Fixed error in the kubernetes list command (a1c3a587)
* Typo and readability fixes (#40) (17007e30)
* Merge pull request #38 from kaihoffman/master (98f6850d)

0.6.45
=============
2020-12-04

* Fixed error handling in all commands (5e789fe3)

0.6.44
=============
2020-12-04

* Update the civogo lib to v0.2.23 (2b8f4d92)
* Merge branch 'master' of https://github.com/civo/cli (9355e877)
* Fixed bug in the install script (585b9d16)

0.6.43
=============
2020-11-29

* Merge pull request #34 from martynbristow/master (af8c5e7d)

0.6.42
=============
2020-11-18

* Updated the civogo lib to v0.2.22 (d5695bff)
* Merge pull request #33 from kaihoffman/master (a7b6ef7e)

0.6.41
=============
2020-11-11

* Added the recycle cmd to kubernetes (d2457212)

0.6.40
=============
2020-11-07

* Fixed the permission in the civo conf file (f89b57e5)

0.6.39
=============
2020-11-02

* Added powershel and fish to the completion cmd (3ce2c969)
* Updated the civogo lib to v0.2.21 (8820f312)

0.6.38
=============
2020-10-16

* Merge pull request #30 from Johannestegner/master (b59581ea)
* Updated the apikey output (f2d000b9)
* Update README.md (318d5988)
* Updated the Change log (a4abc2b0)

0.6.37
=============
2020-10-09

* Add arm64 to build list (52beb72e)

0.6.36
=============
2020-09-25

* Updated the civogo lib to v0.2.19 (b99f4258)

0.6.35
=============
2020-09-23

* Updated the civogo lib to v0.2.18 (9879243a)
* Update README.md (102dd5d6)
* Update the gitignore (b4cce868)
* Update the README.md (999d0895)

0.6.34
=============
2020-09-13

* Allowed SRV record type in the DNS command (594884e1)
* Updated civogo lib to v0.2.17 (10e6d50d)

0.6.33
=============
2020-09-09

* Added STOPPING status (422d3933)

0.6.32
=============
2020-09-08

* improved the custom output (98939932)

0.6.31
=============
2020-09-07

* Fixed error in the custom output (cc05acff)
* Updated the README (c8815d73)
* Update the Change log (3900c15e)

0.6.30
=============
2020-09-04

* Fixed typo error in the Dockerfile (9d479dae)
* Updated Change log (30376d69)
* Fixed error in the Dockerfile (73e5e0d8)
* Update the Change log (6054efff)
* Updated the Dockerfile to add curl (e2d0d5f8)
* Update the Change log (9ed54ddd)
* Added kubectl inside the docker image (65c2c23c)

0.6.29
=============
2020-08-26

* Quota command improved (381ab929)
* Added new show command to apikey (cb8c73d5)

0.6.28
=============
2020-08-24

* Update the civogo lib to v0.2.16 (8257fb46)
* Fixed the config generator (5a2efa5e)
* Remove windows 386 build from goreleaser (ab8d0bba)
* Added verification to all commands (80a158ea)
* Add a little polish for DEVELOPER.md (2c616de4)
* Added the develper file (3fefaffb)

0.6.27
=============
2020-08-20

* Name changed to label in the message of remove (add3b7f4)
* Updated the message in the remove command (c7c57d9d)

0.6.26
=============
2020-08-19

* Update civogo lib to v0.2.15 (09b3bbe2)
* Added verification before delete any object (a9c4334f)

0.6.25
=============
2020-08-19

* Checking correct flag for Kubernetes readiness (e0f543c4)

0.6.24
=============
2020-08-17

* Update on Kubernetes apps documentation (9e0767c7)

0.6.23
=============
2020-08-14

* Added new features when creating a kubernetes cluster (fb51ae8c)

0.6.22
=============
2020-08-14

* Fixed bug in output writing component (53ef3806)

0.6.21
=============
2020-08-11

* Update the civogo to 0.2.14, fixed error in the intance create cmd (78e107de)

0.6.20
=============
2020-08-11

* Fixed error in the instance module and in the color utility (8738797c)

0.6.19
=============
2020-07-30

* Added the auto update option to the cli (eb9ff675)

0.6.18
=============
2020-07-27

* Fixed color in the error message (a30bb996)
* Update Change log (5257cc27)
* Added verification at the moment of delete a snapshot (911cb634)
* Change the lib to add color to the CLI (d35889a5)

0.6.17
=============
2020-07-20

* Fixed error in the configuration of the kubeconfig (87326d25)

0.6.16
=============
2020-07-20

* Merge all color utility in one place (6dd3b564)
* Fixed error in the --merge option for the kubernetes config in windows (1094c066)
* Update the README.md (cfb1259e)

0.6.15
=============
2020-07-18

* Added the --save and --switch option to kubernetes create (950a432d)
* Fixed typo in the kubernetes config error message (bf097440)
* Added verification to kubernetes config cmd (7e4bb504)
* Fix error in kubernetes utils (f1f46612)
* Now if you use --switch with --merge, automatically the cli will change the context in kubernetes (8b9fad39)

0.6.14
=============
2020-07-17

* Removed option egress as direction when creating a firewall rule (7baf88f4)
* Changed the words used to define the direction in the firewall rule (3df9c5b6)
* Added .editorconfig (cfb6ca22)

0.6.13
=============
2020-07-14

* Added verification to the firewall rule creation (a7fc88f9)
* Added new verification step to the kubernetes config (f77cf09e)
* Update the change log (5dc1bcac)

0.6.12
=============
2020-07-08

* Added the option to install multiple application at the same time in the Kubernetes Cluster (611120b1)
* Update the change log (d6cc1d42)

0.6.11
=============
2020-07-07

* Added CPU, RAM and SSD fields to Instance and Kubernetes CMD (6c5cfec7)
* Update the change log (e88016ef)

0.6.10
=============
2020-07-06

* Added new feature (56719687)
* Fixed error in the Makefile (feb35382)
* Merge pull request #8 from ssmiller25/selective-make (6caba09c)
* Chnage the name of the fordel in Homebrew Formula (e96b84a4)
* Merge pull request #7 from ssmiller25/k3s-docs (4a70e48b)
* Fixed error in cli (10196c99)
* Add KubernetesVesion and UpgradeAvailableTo (d61e6a9a)
* Add nfpm and use the v2 of the goreleaser action (49d55069)
* Fixed the message to show after upgrade happened (396934a5)
* Upgrade the version of civogo lib (6ff45d0d)
* Fixed error in the kubernetes cmd (723028c1)
* Change the dockerfile to run as root (7f15c77b)
* fix: Fix error adding initil user (28c99ac9)

0.6.9
=============
2020-07-03

* Chnage the name of the fordel in Homebrew Formula (e96b84a4)
* Merge pull request #7 from ssmiller25/k3s-docs (4a70e48b)

0.6.8
=============
2020-06-30

* Fixed error in cli (10196c99)

0.6.7
=============
2020-06-25

* Add KubernetesVesion and UpgradeAvailableTo (d61e6a9a)

0.6.6
=============
2020-06-24

* Add nfpm and use the v2 of the goreleaser action (49d55069)

0.6.5
=============
2020-06-24

* Fixed the message to show after upgrade happened (396934a5)
* Upgrade the version of civogo lib (6ff45d0d)

0.6.4
=============
2020-06-24

* Fixed error in the kubernetes cmd (723028c1)

0.6.3
=============
2020-06-22

* Change the dockerfile to run as root (7f15c77b)

0.6.2
=============
2020-06-20

* fix: Fix error adding initil user (28c99ac9)
* Update the CHANGELOG (6669b296)
* fix: Fix message in the kubernetes upgrade cmd (05960c0b)
* docs: Update the CHANGELOG (ea04cc57)
* fix: Remove update from the cmd kubernetes show (6cb4144a)
* Added the CHANGELOG.md (e99843a4)
* Added new feature (f1726759)
* Fix correct version showing (83fe22c8)
* Added some feature (a356acaa)
* Fixed some errors (911f046e)
* Remove snap from release process (15d43a7c)
* Update kubernetes config (6413a22a)
* Update .goreleaser.yml (7ff511dd)
* Update .goreleaser.yml (a763c30b)
* Update .goreleaser.yml (3166ee2f)
* Fix error (787ad7de)
* Added some corrections (e5fe1faa)
* Fixed error in goreleaser (72708d8c)
* Modifications in logic (497e0750)
* Tweaks to install.sh (29250021)
* Add configuration to snap (148a0473)
* Update README.md (3bea636c)
* Update the badge (d118e1e3)
* Add badge to the README (0da2f4ff)
* Fix error in the config (8a1430f9)
* Fix user in the Dockerfile (b0d63fe2)
* Update the image to alpine the Dockerfile (0b182a0b)
* Update the Dockerfile (00942fdb)
* Update the cli (4d4ac165)
* Fix wording (0104cacd)
* Update README now this is the main tool (1597d82f)
* Remove comments (f42b14ac)
* Run new test (578d62d7)
* Test brews releaser (425f84f5)
* Try to fix Homebrew releasing (bdc4a966)
* Fix goreleaser.yml (4b0a350b)
* Add homebrew tap updating to goreleaser (2801d339)
* fix(Project): Added some fix (8ed9c0bf)
* Lots of minor tweaks (ef039b50)
* Changes to instance create (82901093)
* feat(Utility): Add TrackTime (7e16e37d)
* feat(Utility): Add random (268021cf)
* fix(Bug): Fix bug found in all code (47c14806)
* fix(Domain): Fix some error (04c0775e)
* feat(Install): Added install bash (2e584464)
* feat(Completion): Added completion cmd (11f44f4b)
* fix(Project): Fix the some bug (a082fdaf)
* fix(Volume): Fix the example (d3e46926)
* refactor(Project): Added example to all cmd (0119e351)
* fix(Network): Remove Name from the table (b5659302)
* fix(LoadBalancer): added new alias (2cb08795)
* refactor(Project): Modify all project error handler (36f932ee)
* fix(Version): fix text in version cmd (d4b80553)
* Merge remote-tracking branch 'origin/master' (732fcc41)
* fix(Version): fix text in version cmd (fe625b9b)
* fix(Version): update the version cmd (b4af168b)
* fix(Version): change the version cmd (17402302)
* fix(Version): update the version cmd (75fe68e8)
* fix(Version): Fixed the version cmd (f01cff84)
* fix(Goreleser): Fix the conf file (3ecf80b0)
* fix(Version): Fixed the version cmd (b35e0cf5)
* fix(Goreleser): Fix the conf file (122998f8)
* fix(Version): Upgrade the version of the civogo lib (c84e2461)
* feat(Version): Added the version cmd (aaa6a091)
* fix(Goreleser): Fix the conf file (f4b78e1b)
* fix(Goreleser): Fix the conf file (6beb0caf)
* fix(Goreleser): Fix the conf file (0953e531)
* fix(Goreleser): Fix the conf file (0ab428ae)
* fix(Goreleser): Fix the action file (2d05ff91)
* fix(Goreleser): Fix the conf file (5b12fdf7)
* fix(Goreleser): Fix the action file (44ad7dff)
* refactor(Project): Modify all project (bc2f508f)
* Merge pull request #2 from civo/dev (daccd7ba)
* Add all but one instance command (7ad58ed4)
* Split more commands out in to their own files (3f3b06e3)
* Split instance commands out to invidual files (2ce29790)
* Add console and reboot commands (d24951af)
* Add github token path to goreleaser (cd811766)
* Add testing workflow (d01e3ed8)
* Add goreleaser support (6410d22b)
* Add the first action command to instances (ec0634b6)
* Change OutputWriter to use labels as well as keys (52eeee28)
* Change to use Find* in civogo (e15de86c)
* Refactor to new OutputWriter (417c47aa)
* Add size command (f7201a09)
* Add quota command (a07f44a6)
* Add Region command (48cd4f0c)
* Complete API key management (bc9e3d53)
* Remove unnecessary comment (7f7fb17f)
* Remove (old) accidentally commited API keys (ea215e4f)
* Complete API key management functionality (94421615)
* Add API Key management (4beb5f1e)
* Update links in README (17848b68)
* Add plan to README (b9fc702c)
* First commit - initial makefile (38300f4e)
* feat(Template): Added the kubernetes cmd (0370dde0)
* feat(Goreleser): Fix the goreleaser file (59a4231c)
* feat(Template): Added the template cmd (4f406582)
* feat(Volume): Added the volume cmd (6a84a030)
* feat(snapshot): Added the snapshot cmd (775885d0)
* feat(network): Added the network cmd (6c6d035e)
* feat(sshkey): Added the ssh key cmd (4f3a6445)
* feat(loadbalancer): Added the cmd loadbalancer (410abb51)
* feat(firewall): Added the cmd firewall, and firewall rule (0ad42aed)
* feat(domain): Added the cmd domain record (51de188d)
* feat(domain): Added the cmd domain (698df535)
* feat(instance): Added the option to create instances (22dcec05)
* refactor(apikey): Modify the apikey (1d113038)
* refactor(apikey): Modify the apikey (430a4114)
* refactor(apikey): Modify the apikey function (a8310907)
* Update app name in tap-release (de91a6e0)
* Create automatic tap release (fa017dff)

0.6.1
=============
2020-06-19

* Update the CHANGELOG (6669b296)
* fix: Fix message in the kubernetes upgrade cmd (05960c0b)
* docs: Update the CHANGELOG (ea04cc57)
* fix: Remove update from the cmd kubernetes show (6cb4144a)
* Added the CHANGELOG.md (e99843a4)
* Added new feature (f1726759)
* Fix correct version showing (83fe22c8)

0.6.0
=============
2020-06-18

* Added some feature (a356acaa)

0.2.3
=============
2020-06-18

* Fixed some errors (911f046e)
* Remove snap from release process (15d43a7c)
* Update kubernetes config (6413a22a)

0.2.2
=============
2020-06-17

* Update .goreleaser.yml (7ff511dd)
* Update .goreleaser.yml (a763c30b)
* Update .goreleaser.yml (3166ee2f)
* Fix error (787ad7de)
* Added some corrections (e5fe1faa)

0.2.1
=============
2020-06-17

* Fixed error in goreleaser (72708d8c)

0.2.0
=============
2020-06-17

* Modifications in logic (497e0750)
* Tweaks to install.sh (29250021)
* Add configuration to snap (148a0473)
* Update README.md (3bea636c)
* Update the badge (d118e1e3)
* Add badge to the README (0da2f4ff)
* Fix error in the config (8a1430f9)
* Fix user in the Dockerfile (b0d63fe2)
* Update the image to alpine the Dockerfile (0b182a0b)
* Update the Dockerfile (00942fdb)
* Update the cli (4d4ac165)
* Fix wording (0104cacd)
* Update README now this is the main tool (1597d82f)
* Remove comments (f42b14ac)
* Run new test (578d62d7)
* Test brews releaser (425f84f5)
* Try to fix Homebrew releasing (bdc4a966)
* Fix goreleaser.yml (4b0a350b)
* Add homebrew tap updating to goreleaser (2801d339)
* fix(Project): Added some fix (8ed9c0bf)
* Lots of minor tweaks (ef039b50)
* Changes to instance create (82901093)
* feat(Utility): Add TrackTime (7e16e37d)
* feat(Utility): Add random (268021cf)
* fix(Bug): Fix bug found in all code (47c14806)

0.1.19
=============
2020-06-17

* Tweaks to install.sh (29250021)
* Add configuration to snap (148a0473)

0.1.18
=============
2020-06-16

* Update README.md (3bea636c)
* Update the badge (d118e1e3)
* Add badge to the README (0da2f4ff)
* Fix error in the config (8a1430f9)
* Fix user in the Dockerfile (b0d63fe2)

0.1.17
=============
2020-06-16

* Update the image to alpine the Dockerfile (0b182a0b)

0.1.16
=============
2020-06-16

* Update the Dockerfile (00942fdb)
* Update the cli (4d4ac165)
* Fix wording (0104cacd)
* Update README now this is the main tool (1597d82f)
* Remove comments (f42b14ac)

0.1.15
=============
2020-06-15

* Run new test (578d62d7)
* Test brews releaser (425f84f5)

0.1.14
=============
2020-06-12

* Try to fix Homebrew releasing (bdc4a966)

0.1.13
=============
2020-06-12

* Fix goreleaser.yml (4b0a350b)
* Add homebrew tap updating to goreleaser (2801d339)

0.1.12
=============
2020-05-21

* fix(Project): Added some fix (8ed9c0bf)

0.1.11
=============
2020-05-20

* Lots of minor tweaks (ef039b50)
* Changes to instance create (82901093)

0.1.10
=============
2020-05-19

* feat(Utility): Add TrackTime (7e16e37d)
* feat(Utility): Add random (268021cf)
* fix(Bug): Fix bug found in all code (47c14806)
* fix(Domain): Fix some error (04c0775e)
* feat(Install): Added install bash (2e584464)
* feat(Completion): Added completion cmd (11f44f4b)
* fix(Project): Fix the some bug (a082fdaf)
* fix(Volume): Fix the example (d3e46926)
* refactor(Project): Added example to all cmd (0119e351)
* fix(Network): Remove Name from the table (b5659302)
* fix(LoadBalancer): added new alias (2cb08795)
* refactor(Project): Modify all project error handler (36f932ee)
* fix(Version): fix text in version cmd (d4b80553)
* Merge remote-tracking branch 'origin/master' (732fcc41)
* fix(Version): fix text in version cmd (fe625b9b)
* fix(Version): update the version cmd (b4af168b)
* fix(Version): change the version cmd (17402302)
* fix(Version): update the version cmd (75fe68e8)
* fix(Version): Fixed the version cmd (f01cff84)
* fix(Goreleser): Fix the conf file (3ecf80b0)
* fix(Version): Fixed the version cmd (b35e0cf5)
* fix(Goreleser): Fix the conf file (122998f8)
* fix(Version): Upgrade the version of the civogo lib (c84e2461)
* feat(Version): Added the version cmd (aaa6a091)
* fix(Goreleser): Fix the conf file (f4b78e1b)
* fix(Goreleser): Fix the conf file (6beb0caf)
* fix(Goreleser): Fix the conf file (0953e531)
* fix(Goreleser): Fix the conf file (0ab428ae)
* fix(Goreleser): Fix the action file (2d05ff91)
* fix(Goreleser): Fix the conf file (5b12fdf7)
* fix(Goreleser): Fix the action file (44ad7dff)
* refactor(Project): Modify all project (bc2f508f)
* Merge pull request #2 from civo/dev (daccd7ba)
* Add all but one instance command (7ad58ed4)
* Split more commands out in to their own files (3f3b06e3)
* Split instance commands out to invidual files (2ce29790)
* Add console and reboot commands (d24951af)
* Add github token path to goreleaser (cd811766)
* Add testing workflow (d01e3ed8)
* Add goreleaser support (6410d22b)
* Add the first action command to instances (ec0634b6)
* Change OutputWriter to use labels as well as keys (52eeee28)
* Change to use Find* in civogo (e15de86c)
* Refactor to new OutputWriter (417c47aa)
* Add size command (f7201a09)
* Add quota command (a07f44a6)
* Add Region command (48cd4f0c)
* Complete API key management (bc9e3d53)
* Remove unnecessary comment (7f7fb17f)
* Remove (old) accidentally commited API keys (ea215e4f)
* Complete API key management functionality (94421615)
* Add API Key management (4beb5f1e)
* Update links in README (17848b68)
* Add plan to README (b9fc702c)
* First commit - initial makefile (38300f4e)
* feat(Template): Added the kubernetes cmd (0370dde0)
* feat(Goreleser): Fix the goreleaser file (59a4231c)
* feat(Template): Added the template cmd (4f406582)
* feat(Volume): Added the volume cmd (6a84a030)
* feat(snapshot): Added the snapshot cmd (775885d0)
* feat(network): Added the network cmd (6c6d035e)
* feat(sshkey): Added the ssh key cmd (4f3a6445)
* feat(loadbalancer): Added the cmd loadbalancer (410abb51)
* feat(firewall): Added the cmd firewall, and firewall rule (0ad42aed)
* feat(domain): Added the cmd domain record (51de188d)
* feat(domain): Added the cmd domain (698df535)
* feat(instance): Added the option to create instances (22dcec05)
* refactor(apikey): Modify the apikey (1d113038)
* refactor(apikey): Modify the apikey (430a4114)
* refactor(apikey): Modify the apikey function (a8310907)
* Update app name in tap-release (de91a6e0)
* Create automatic tap release (fa017dff)

0.1.9
=============
2020-05-18

* fix(Domain): Fix some error (04c0775e)
* feat(Install): Added install bash (2e584464)

0.1.8
=============
2020-05-14

* feat(Completion): Added completion cmd (11f44f4b)

0.1.7
=============
2020-05-14

* fix(Project): Fix the some bug (a082fdaf)
* fix(Volume): Fix the example (d3e46926)
* refactor(Project): Added example to all cmd (0119e351)
* fix(Network): Remove Name from the table (b5659302)
* fix(LoadBalancer): added new alias (2cb08795)
* refactor(Project): Modify all project error handler (36f932ee)
* fix(Version): fix text in version cmd (d4b80553)

0.1.6
=============
2020-05-13

* Merge remote-tracking branch 'origin/master' (732fcc41)
* fix(Version): fix text in version cmd (fe625b9b)

0.1.5
=============
2020-05-13

* fix(Version): update the version cmd (b4af168b)
* fix(Version): change the version cmd (17402302)
* fix(Version): update the version cmd (75fe68e8)

0.1.4
=============
2020-05-13

* fix(Version): Fixed the version cmd (f01cff84)

0.1.3
=============
2020-05-13

* fix(Goreleser): Fix the conf file (3ecf80b0)
* fix(Version): Fixed the version cmd (b35e0cf5)
* fix(Goreleser): Fix the conf file (122998f8)
* fix(Version): Upgrade the version of the civogo lib (c84e2461)
* feat(Version): Added the version cmd (aaa6a091)

0.1.2
=============
2020-05-11

* fix(Goreleser): Fix the conf file (f4b78e1b)
* fix(Goreleser): Fix the conf file (6beb0caf)
* fix(Goreleser): Fix the conf file (0953e531)
* fix(Goreleser): Fix the conf file (0ab428ae)
* fix(Goreleser): Fix the action file (2d05ff91)
* fix(Goreleser): Fix the conf file (5b12fdf7)
* fix(Goreleser): Fix the action file (44ad7dff)
* refactor(Project): Modify all project (bc2f508f)
* Merge pull request #2 from civo/dev (daccd7ba)
* Add all but one instance command (7ad58ed4)
* Split more commands out in to their own files (3f3b06e3)
* Split instance commands out to invidual files (2ce29790)
* Add console and reboot commands (d24951af)
* Add github token path to goreleaser (cd811766)
* Add testing workflow (d01e3ed8)
* Add goreleaser support (6410d22b)
* Add the first action command to instances (ec0634b6)
* Change OutputWriter to use labels as well as keys (52eeee28)
* Change to use Find* in civogo (e15de86c)
* Refactor to new OutputWriter (417c47aa)
* Add size command (f7201a09)
* Add quota command (a07f44a6)
* Add Region command (48cd4f0c)
* Complete API key management (bc9e3d53)
* Remove unnecessary comment (7f7fb17f)
* Remove (old) accidentally commited API keys (ea215e4f)
* Complete API key management functionality (94421615)
* Add API Key management (4beb5f1e)
* Update links in README (17848b68)
* Add plan to README (b9fc702c)
* First commit - initial makefile (38300f4e)
* feat(Template): Added the kubernetes cmd (0370dde0)
* feat(Goreleser): Fix the goreleaser file (59a4231c)
* feat(Template): Added the template cmd (4f406582)
* feat(Volume): Added the volume cmd (6a84a030)
* feat(snapshot): Added the snapshot cmd (775885d0)
* feat(network): Added the network cmd (6c6d035e)
* feat(sshkey): Added the ssh key cmd (4f3a6445)
* feat(loadbalancer): Added the cmd loadbalancer (410abb51)
* feat(firewall): Added the cmd firewall, and firewall rule (0ad42aed)
* feat(domain): Added the cmd domain record (51de188d)
* feat(domain): Added the cmd domain (698df535)
* feat(instance): Added the option to create instances (22dcec05)
* refactor(apikey): Modify the apikey (1d113038)
* refactor(apikey): Modify the apikey (430a4114)
* refactor(apikey): Modify the apikey function (a8310907)

0.1.1
=============
2020-03-04

* Add github token path to goreleaser (cd811766)
* Add testing workflow (d01e3ed8)

0.1.0
=============
2020-03-04

* Add goreleaser support (6410d22b)
* Add the first action command to instances (ec0634b6)
* Change OutputWriter to use labels as well as keys (52eeee28)
* Change to use Find* in civogo (e15de86c)
* Refactor to new OutputWriter (417c47aa)
* Add size command (f7201a09)
* Add quota command (a07f44a6)
* Add Region command (48cd4f0c)
* Complete API key management (bc9e3d53)
* Remove unnecessary comment (7f7fb17f)
* Remove (old) accidentally commited API keys (ea215e4f)
* Complete API key management functionality (94421615)
* Add API Key management (4beb5f1e)
* Update links in README (17848b68)
* Add plan to README (b9fc702c)
* First commit - initial makefile (38300f4e)


