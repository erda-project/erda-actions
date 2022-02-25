package main

// func TestGetSonarIssues(t *testing.T) {
// 	body, err := getSonarIssues("https://sonar.app.terminus.io",
// 		"xxx",
// 		"repo.endpoints.showcase-front_6a84e35432851eedbcaba73f28d970687fb72b7d",
// 		"VULNERABILITY",
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(err)

// 	fmt.Println(string(body))
// }

// func TestGetSonarIssuesTree(t *testing.T) {
// 	body, err := getSonarIssuesTree("https://sonar.app.terminus.io",
// 		"xxx",
// 		"js-co",
// 		"duplications",
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(err)

// 	fmt.Println(string(body))
// }

// func TestSonar_GetIssuesStatistics(t *testing.T) {
// 	sonar := Sonar{}
// 	body, err := sonar.GetIssuesStatistics("repo.endpoints.showcase-front_6a84e35432851eedbcaba73f28d970687fb72b7d",
// 		"",
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(err)

// 	fmt.Printf("body:%v", body)
// }

// func TestSonar_GetIssues(t *testing.T) {
// 	sonar := Sonar{}
// 	body, err := sonar.GetIssues(nil, "https://sonar.app.terminus.io",
// 		"repo.endpoints.showcase-front_6a84e35432851eedbcaba73f28d970687fb72b7d",
// 		"BUG",
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(err)

// 	fmt.Printf("body:%v", body)
// }

// func TestSonar_GetIssuesTree(t *testing.T) {
// 	sonar := Sonar{}
// 	body, err := sonar.GetIssuesTree("https://sonar.app.terminus.io",
// 		"xxx",
// 		"repo.endpoints.showcase-front_6a84e35432851eedbcaba73f28d970687fb72b7d",
// 		"coverage", "",
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Printf("body:%v", body)
// }

// func TestSonar_ComposeIssues(t *testing.T) {
// 	sonar := Sonar{}
// 	store := sonar.ComposeIssues("https://sonar.app.terminus.io",
// 		"repo.services.pmp-backend_.tmp.b9afcac9054676524a3b61e9a67be14cec89f84a",
// 		nil,
// 	)

// 	fmt.Printf("+++store:%v\n", store)
// }

// func TestDelSonarServerProject(t *testing.T) {
// 	delSonarServerProject("http://sonar.app.terminus.io",
// 		"repo.services.showcase_b0320423546d702cd3aa4019396de445672d31ef",
// 		"",
// 	)
// }

// func TestGetIssuesPath(t *testing.T) {
// 	fmt.Println(getIssuesPath("sonar-resource/sonar", "sonar_test.go"))
// 	fmt.Println(getIssuesPath("sonar", "sonar_test.go"))
// }

// func TestMakeSonarPath(t *testing.T) {
// 	sonarPath := makeSonarPath("pandora-search/src/main/java/io/terminus/pandora/search/service/ESNodeService.java",
// 		"",
// 		"repo")
// 	t.Log(sonarPath)
// 	assert.Equal(t, sonarPath, "pandora-search/src/main/java/io/terminus/pandora/search/service/ESNodeService.java")

// 	sonarPath = makeSonarPath("",
// 		"src/main/java/io/terminus/pandora/search/service/ESNodeService.java",
// 		"repo")
// 	t.Log(sonarPath)
// 	assert.Equal(t, sonarPath, "src/main/java/io/terminus/pandora/search/service/ESNodeService.java")

// 	sonarPath = makeSonarPath("blog-web/src/main/java/io/terminus/pampas/showcase/Application.java",
// 		"",
// 		"repo/services/showcase")
// 	t.Log(sonarPath)
// 	assert.Equal(t, sonarPath, "services/showcase/blog-web/src/main/java/io/terminus/pampas/showcase/Application.java")

// 	sonarPath = makeSonarPath("",
// 		"src/main/java/io/terminus/pampas/showcase/Application.java",
// 		"repo/services/showcase")
// 	t.Log(sonarPath)
// 	assert.Equal(t, sonarPath, "services/showcase/src/main/java/io/terminus/pampas/showcase/Application.java")
// }
