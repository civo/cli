
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


