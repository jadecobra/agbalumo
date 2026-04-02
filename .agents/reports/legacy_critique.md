# ChiefCritic Technical Debt Report
Generated: Wed Apr  1 11:47:31 CDT 2026

## 1. Cognitive Complexity (Threshold < 10)
```
91 agent ExtractRoutes internal/agent/ast.go:21:1
79 agent_test TestStateSerialization internal/agent/state_test.go:13:1
78 commands VerifyCmd cmd/harness/commands/verify.go:13:1
66 agent ExtractTemplateFunctionCalls internal/agent/verify.go:432:1
65 commands ChaosCmd cmd/harness/commands/chaos.go:21:1
54 service (*CSVService).ParseAndImport internal/service/csv.go:25:1
51 service TestParseAndImport internal/service/csv_test.go:11:1
49 agent ExtractCLIMarkdownCommands internal/agent/drift.go:171:1
42 agent VerifySecurityStatic internal/agent/security.go:70:1
41 agent VerifyRedTest internal/agent/verify.go:57:1
39 agent checkXSS internal/agent/security.go:341:1
36 agent checkStructuralRaw internal/agent/security.go:563:1
32 service (*LocalImageService).UploadImage internal/service/image.go:48:1
29 agent checkSSRF internal/agent/security.go:396:1
28 agent ParseTestJSON internal/agent/redtest.go:31:1
23 sqlite_test BenchmarkSearchPerformance internal/repository/sqlite/search_performance_test.go:15:1
22 agent ParseMarkdownTracker internal/agent/progress.go:32:1
22 domain (*Listing).Validate internal/domain/listing.go:121:1
22 agent ExtractCLICodeCommands internal/agent/drift.go:136:1
21 sqlite_test TestSaveFeedback internal/repository/sqlite/feedback_test.go:16:1
21 agent VerifyCoverage internal/agent/verify.go:304:1
20 agent CalculateContextCost internal/agent/cost.go:24:1
20 agent checkFileInclusion internal/agent/security.go:448:1
20 listing (*ListingHandler).HandleUpdate internal/module/listing/listing_mutations.go:54:1
19 ui BuildGlobalFuncMap internal/ui/renderer.go:62:1
19 agent ParseCoverageProfile internal/agent/coverage.go:18:1
19 sqlite (*SQLiteRepository).BulkInsertListings internal/repository/sqlite/sqlite_listing_write.go:121:1
18 cmd TestCLIJSONOutput cmd/cli_json_test.go:14:1
18 middleware_test TestSecureHeaders internal/middleware/security_test.go:13:1
18 agent TestVerifyRedTest internal/agent/verify_test.go:302:1
18 admin (*AdminHandler).HandleBulkAction internal/module/admin/admin_bulk.go:14:1
18 commands GateCmd cmd/harness/commands/gate.go:11:1
18 service (*GoogleGeocodingService).GetCity internal/service/geocoding.go:34:1
17 agent TestCalculateContextCost internal/agent/cost_test.go:9:1
17 agent checkInsecurePatternsGo internal/agent/security.go:527:1
16 listing_test TestTemplateTailwindCleanup internal/module/listing/ui_regression_home_test.go:62:1
16 listing_test TestListingHandler_FormParsing internal/module/listing/listing_form_integration_test.go:17:1
16 sqlite (*SQLiteRepository).FindAll internal/repository/sqlite/sqlite_listing_read.go:54:1
15 agent checkSQLi internal/agent/security.go:302:1
15 listing (*ListingHandler).processAndSave internal/module/listing/listing_mutations.go:136:1
14 commands UpdateCoverageCmd cmd/harness/commands/update_coverage.go:12:1
14 ui TestTemplateRenderer_Render_IsNew internal/ui/renderer_test.go:94:1
14 agent VerifyApiSpec internal/agent/verify.go:156:1
14 admin (*AdminHandler).HandleToggleFeatured internal/module/admin/admin_listings.go:74:1
14 agent TestVerifyApiSpec_ExtractMarkdownRoutes_Direct internal/agent/verify_apispec_test.go:608:1
14 agent TestVerifyLint internal/agent/verify_test.go:85:1
14 cmd ResolveServerConfig cmd/serve.go:21:1
14 commands CostCmd cmd/harness/commands/cost.go:10:1
14 admin_test TestAdminHandler_HandleToggleFeatured internal/module/admin/admin_actions_test.go:31:1
13 agent ArchivePassedCategories internal/agent/progress.go:101:1
13 commands TestRootFunctions cmd/harness/commands/root_test.go:8:1
13 sqlite_test TestSaveCategory_PreservesExistingCategories internal/repository/sqlite/sqlite_category_test.go:260:1
13 service (*CSVService).parseRow internal/service/csv.go:226:1
13 admin (*AdminHandler).HandleAddCategory internal/module/admin/admin.go:232:1
13 ui compileTemplates internal/ui/renderer.go:128:1
13 cmd_test TestUserRoutes cmd/server_user_test.go:14:1
13 seeder GenerateStressListings internal/seeder/stress_generator.go:21:1
13 commands InitCmd cmd/harness/commands/init.go:10:1
13 seeder_test TestEnsureCategoriesSeeded_Verification internal/seeder/config_verification_test.go:16:1
13 history TestStore internal/history/history_test.go:10:1
12 commands checkAndApplyProgressUpdate cmd/harness/commands/root.go:99:1
12 commands summarizeProgress cmd/harness/commands/root.go:59:1
12 sqlite_test TestGetAllFeedback internal/repository/sqlite/feedback_test.go:103:1
12 agent isIgnored internal/agent/security.go:506:1
12 agent TestStateExtraCoverage internal/agent/state_extra_test.go:8:1
12 agent TestVerifySecurityStaticGate internal/agent/verify_security_test.go:10:1
12 admin (*AdminHandler).HandleBulkUpload internal/module/admin/admin_bulk.go:80:1
12 cmd TestListingCommands_RunCoverage cmd/listing_cmd_test.go:14:1
12 main run cmd/aglog/main.go:43:1
12 admin (*AdminHandler).HandleDashboard internal/module/admin/admin.go:155:1
12 handler_test TestIntegration_DataValidation internal/handler/integration_ds_test.go:19:1
12 commands HandoffCmd cmd/harness/commands/handoff.go:11:1
12 commands TestChaosCommand cmd/harness/commands/chaos_test.go:14:1
12 agent TestVerifyImplementation internal/agent/verify_test.go:34:1
11 main runAudit cmd/security-audit/main.go:90:1
11 listing parseEventDates internal/module/listing/listing_form.go:93:1
11 main TestCheckVuln cmd/security-audit/audit_test.go:159:1
11 agent VerifyImplementation internal/agent/verify.go:236:1
11 ui TestTemplateRenderer_EdgeCases internal/ui/renderer_test.go:276:1
11 util TestUniqueStrings internal/util/fs_test.go:201:1
11 main TestCheckHeaders cmd/security-audit/audit_test.go:12:1
11 agent checkEntropyGo internal/agent/security.go:237:1
11 admin_test TestAdminHandler_HandleLoginAction internal/module/admin/admin_login_test.go:62:1
11 agent TestTaskfileToolingOptimization internal/agent/task_optimization_test.go:9:1
10 main checkHeaders cmd/security-audit/main.go:162:1
10 agent EnforceCoverage internal/agent/coverage.go:146:1
10 sqlite_test TestReproCategoryRegression internal/repository/sqlite/repro_test.go:13:1
10 admin_test TestAdminHandler_RegisterRoutes internal/module/admin/admin_routes_test.go:26:1
10 service (*CSVService).GenerateCSV internal/service/csv.go:142:1
10 commands SetPhaseCmd cmd/harness/commands/set_phase.go:9:1
```
## 2. Repeated Constants
```
goconst: 2026/04/01 11:47:31 Found 14 Go files to process in batches of 50
goconst: 2026/04/01 11:47:31 Processing batch 1/1 (14 files)
cmd/stress.go:31:46:2 other occurrence(s) of "count" found in: cmd/stress.go:36:38 cmd/stress.go:48:42
cmd/stress.go:36:38:2 other occurrence(s) of "count" found in: cmd/stress.go:31:46 cmd/stress.go:48:42
cmd/stress.go:48:42:2 other occurrence(s) of "count" found in: cmd/stress.go:31:46 cmd/stress.go:36:38
cmd/listing.go:61:49:2 other occurrence(s) of "city" found in: cmd/listing_backfill.go:58:60 cmd/listing.go:80:49
cmd/listing_backfill.go:58:60:2 other occurrence(s) of "city" found in: cmd/listing.go:61:49 cmd/listing.go:80:49
cmd/listing.go:80:49:2 other occurrence(s) of "city" found in: cmd/listing.go:61:49 cmd/listing_backfill.go:58:60
cmd/listing_update.go:22:15:3 other occurrence(s) of "Listing not found" found in: cmd/admin.go:31:15 cmd/admin.go:54:15 cmd/admin.go:77:15
cmd/admin.go:31:15:3 other occurrence(s) of "Listing not found" found in: cmd/listing_update.go:22:15 cmd/admin.go:54:15 cmd/admin.go:77:15
cmd/admin.go:54:15:3 other occurrence(s) of "Listing not found" found in: cmd/listing_update.go:22:15 cmd/admin.go:31:15 cmd/admin.go:77:15
cmd/admin.go:77:15:3 other occurrence(s) of "Listing not found" found in: cmd/listing_update.go:22:15 cmd/admin.go:31:15 cmd/admin.go:54:15
cmd/listing_update.go:90:15:1 other occurrence(s) of "Failed to update listing" found in: cmd/admin.go:83:15
cmd/admin.go:83:15:1 other occurrence(s) of "Failed to update listing" found in: cmd/listing_update.go:90:15
cmd/category.go:33:39:1 other occurrence(s) of "claimable" found in: cmd/category.go:89:31
cmd/category.go:89:31:1 other occurrence(s) of "claimable" found in: cmd/category.go:33:39
cmd/listing.go:58:62:1 other occurrence(s) of "Business" found in: cmd/benchmark.go:48:37
cmd/benchmark.go:48:37:1 other occurrence(s) of "Business" found in: cmd/listing.go:58:62
cmd/listing.go:75:51:1 other occurrence(s) of "company" found in: cmd/listing.go:94:51
cmd/listing.go:94:51:1 other occurrence(s) of "company" found in: cmd/listing.go:75:51
cmd/serve.go:95:48:1 other occurrence(s) of "addr" found in: cmd/serve.go:104:54
cmd/serve.go:104:54:1 other occurrence(s) of "addr" found in: cmd/serve.go:95:48
cmd/listing.go:69:52:1 other occurrence(s) of "deadline" found in: cmd/listing.go:88:52
cmd/listing.go:88:52:1 other occurrence(s) of "deadline" found in: cmd/listing.go:69:52
cmd/listing.go:60:56:1 other occurrence(s) of "description" found in: cmd/listing.go:79:56
cmd/listing.go:79:56:1 other occurrence(s) of "description" found in: cmd/listing.go:60:56
cmd/listing.go:57:50:2 other occurrence(s) of "title" found in: cmd/listing.go:78:50 cmd/listing.go:97:40
cmd/listing.go:78:50:2 other occurrence(s) of "title" found in: cmd/listing.go:57:50 cmd/listing.go:97:40
cmd/listing.go:97:40:2 other occurrence(s) of "title" found in: cmd/listing.go:57:50 cmd/listing.go:78:50
cmd/listing.go:62:52:4 other occurrence(s) of "address" found in: cmd/listing_backfill.go:42:59 cmd/listing_backfill.go:45:58 cmd/listing_backfill.go:60:65 cmd/listing.go:81:52
cmd/listing_backfill.go:42:59:4 other occurrence(s) of "address" found in: cmd/listing.go:62:52 cmd/listing_backfill.go:45:58 cmd/listing_backfill.go:60:65 cmd/listing.go:81:52
cmd/listing_backfill.go:45:58:4 other occurrence(s) of "address" found in: cmd/listing.go:62:52 cmd/listing_backfill.go:42:59 cmd/listing_backfill.go:60:65 cmd/listing.go:81:52
cmd/listing_backfill.go:60:65:4 other occurrence(s) of "address" found in: cmd/listing.go:62:52 cmd/listing_backfill.go:42:59 cmd/listing_backfill.go:45:58 cmd/listing.go:81:52
cmd/listing.go:81:52:4 other occurrence(s) of "address" found in: cmd/listing.go:62:52 cmd/listing_backfill.go:42:59 cmd/listing_backfill.go:45:58 cmd/listing_backfill.go:60:65
cmd/listing.go:67:53:1 other occurrence(s) of "image-url" found in: cmd/listing.go:86:53
cmd/listing.go:86:53:1 other occurrence(s) of "image-url" found in: cmd/listing.go:67:53
cmd/listing.go:72:50:1 other occurrence(s) of "skills" found in: cmd/listing.go:91:50
cmd/listing.go:91:50:1 other occurrence(s) of "skills" found in: cmd/listing.go:72:50
cmd/listing.go:74:52:1 other occurrence(s) of "apply-url" found in: cmd/listing.go:93:52
cmd/listing.go:93:52:1 other occurrence(s) of "apply-url" found in: cmd/listing.go:74:52
cmd/listing.go:64:50:1 other occurrence(s) of "phone" found in: cmd/listing.go:83:50
cmd/listing.go:83:50:1 other occurrence(s) of "phone" found in: cmd/listing.go:64:50
cmd/listing_read.go:13:9:1 other occurrence(s) of "list" found in: cmd/category.go:53:9
cmd/category.go:53:9:1 other occurrence(s) of "list" found in: cmd/listing_read.go:13:9
cmd/listing_read.go:27:42:29 other occurrence(s) of "error" found in: cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/listing_read.go:62:40:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/stress.go:25:52:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/stress.go:38:49:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/seed.go:21:52:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/listing_delete.go:20:43:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/benchmark.go:25:52:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/serve.go:78:41:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/listing_create.go:55:63:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/listing_update.go:22:36:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/listing_update.go:90:43:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/listing_create.go:75:43:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/category.go:44:42:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/admin.go:31:36:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/admin.go:37:44:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/admin.go:54:36:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/category.go:60:43:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/listing_backfill.go:33:43:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/listing_backfill.go:45:80:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/listing_backfill.go:53:74:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/benchmark.go:63:52:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/admin.go:60:43:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/admin.go:77:36:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/admin.go:83:43:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/server.go:252:43:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/admin.go:103:47:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/admin.go:138:38:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:177:33 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/admin.go:177:33:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:183:41 cmd/listing.go:104:47
cmd/admin.go:183:41:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/listing.go:104:47
cmd/listing.go:104:47:29 other occurrence(s) of "error" found in: cmd/listing_read.go:27:42 cmd/listing_read.go:62:40 cmd/stress.go:25:52 cmd/stress.go:38:49 cmd/seed.go:21:52 cmd/listing_delete.go:20:43 cmd/benchmark.go:25:52 cmd/serve.go:78:41 cmd/listing_create.go:55:63 cmd/listing_update.go:22:36 cmd/listing_update.go:90:43 cmd/listing_create.go:75:43 cmd/category.go:44:42 cmd/admin.go:31:36 cmd/admin.go:37:44 cmd/admin.go:54:36 cmd/category.go:60:43 cmd/listing_backfill.go:33:43 cmd/listing_backfill.go:45:80 cmd/listing_backfill.go:53:74 cmd/benchmark.go:63:52 cmd/admin.go:60:43 cmd/admin.go:77:36 cmd/admin.go:83:43 cmd/server.go:252:43 cmd/admin.go:103:47 cmd/admin.go:138:38 cmd/admin.go:177:33 cmd/admin.go:183:41
cmd/listing_create.go:52:28:3 other occurrence(s) of "2006-01-02" found in: cmd/listing_update.go:57:28 cmd/admin.go:125:76 cmd/listing.go:140:57
cmd/listing_update.go:57:28:3 other occurrence(s) of "2006-01-02" found in: cmd/listing_create.go:52:28 cmd/admin.go:125:76 cmd/listing.go:140:57
cmd/admin.go:125:76:3 other occurrence(s) of "2006-01-02" found in: cmd/listing_create.go:52:28 cmd/listing_update.go:57:28 cmd/listing.go:140:57
cmd/listing.go:140:57:3 other occurrence(s) of "2006-01-02" found in: cmd/listing_create.go:52:28 cmd/listing_update.go:57:28 cmd/admin.go:125:76
cmd/stress.go:25:36:2 other occurrence(s) of "path" found in: cmd/seed.go:21:36 cmd/benchmark.go:25:36
cmd/seed.go:21:36:2 other occurrence(s) of "path" found in: cmd/stress.go:25:36 cmd/benchmark.go:25:36
cmd/benchmark.go:25:36:2 other occurrence(s) of "path" found in: cmd/stress.go:25:36 cmd/seed.go:21:36
cmd/listing.go:65:53:1 other occurrence(s) of "whatsapp" found in: cmd/listing.go:84:53
cmd/listing.go:84:53:1 other occurrence(s) of "whatsapp" found in: cmd/listing.go:65:53
cmd/listing_update.go:62:28:5 other occurrence(s) of "2006-01-02T15:04" found in: cmd/listing_update.go:67:28 cmd/listing_update.go:75:28 cmd/listing_create.go:59:28 cmd/listing_create.go:64:28 cmd/listing_create.go:69:28
cmd/listing_update.go:67:28:5 other occurrence(s) of "2006-01-02T15:04" found in: cmd/listing_update.go:62:28 cmd/listing_update.go:75:28 cmd/listing_create.go:59:28 cmd/listing_create.go:64:28 cmd/listing_create.go:69:28
cmd/listing_update.go:75:28:5 other occurrence(s) of "2006-01-02T15:04" found in: cmd/listing_update.go:62:28 cmd/listing_update.go:67:28 cmd/listing_create.go:59:28 cmd/listing_create.go:64:28 cmd/listing_create.go:69:28
cmd/listing_create.go:59:28:5 other occurrence(s) of "2006-01-02T15:04" found in: cmd/listing_update.go:62:28 cmd/listing_update.go:67:28 cmd/listing_update.go:75:28 cmd/listing_create.go:64:28 cmd/listing_create.go:69:28
cmd/listing_create.go:64:28:5 other occurrence(s) of "2006-01-02T15:04" found in: cmd/listing_update.go:62:28 cmd/listing_update.go:67:28 cmd/listing_update.go:75:28 cmd/listing_create.go:59:28 cmd/listing_create.go:69:28
cmd/listing_create.go:69:28:5 other occurrence(s) of "2006-01-02T15:04" found in: cmd/listing_update.go:62:28 cmd/listing_update.go:67:28 cmd/listing_update.go:75:28 cmd/listing_create.go:59:28 cmd/listing_create.go:64:28
cmd/listing.go:71:52:1 other occurrence(s) of "event-end" found in: cmd/listing.go:90:52
cmd/listing.go:90:52:1 other occurrence(s) of "event-end" found in: cmd/listing.go:71:52
cmd/serve.go:39:25:6 other occurrence(s) of "production" found in: cmd/serve.go:47:12 cmd/server.go:38:16 cmd/server.go:96:25 cmd/server.go:109:57 cmd/serve.go:101:14 cmd/server.go:255:16
cmd/serve.go:47:12:6 other occurrence(s) of "production" found in: cmd/serve.go:39:25 cmd/server.go:38:16 cmd/server.go:96:25 cmd/server.go:109:57 cmd/serve.go:101:14 cmd/server.go:255:16
cmd/server.go:38:16:6 other occurrence(s) of "production" found in: cmd/serve.go:39:25 cmd/serve.go:47:12 cmd/server.go:96:25 cmd/server.go:109:57 cmd/serve.go:101:14 cmd/server.go:255:16
cmd/server.go:96:25:6 other occurrence(s) of "production" found in: cmd/serve.go:39:25 cmd/serve.go:47:12 cmd/server.go:38:16 cmd/server.go:109:57 cmd/serve.go:101:14 cmd/server.go:255:16
cmd/server.go:109:57:6 other occurrence(s) of "production" found in: cmd/serve.go:39:25 cmd/serve.go:47:12 cmd/server.go:38:16 cmd/server.go:96:25 cmd/serve.go:101:14 cmd/server.go:255:16
cmd/serve.go:101:14:6 other occurrence(s) of "production" found in: cmd/serve.go:39:25 cmd/serve.go:47:12 cmd/server.go:38:16 cmd/server.go:96:25 cmd/server.go:109:57 cmd/server.go:255:16
cmd/server.go:255:16:6 other occurrence(s) of "production" found in: cmd/serve.go:39:25 cmd/serve.go:47:12 cmd/server.go:38:16 cmd/server.go:96:25 cmd/server.go:109:57 cmd/serve.go:101:14
cmd/stress.go:25:15:2 other occurrence(s) of "Failed to open DB" found in: cmd/seed.go:21:15 cmd/benchmark.go:25:15
cmd/seed.go:21:15:2 other occurrence(s) of "Failed to open DB" found in: cmd/stress.go:25:15 cmd/benchmark.go:25:15
cmd/benchmark.go:25:15:2 other occurrence(s) of "Failed to open DB" found in: cmd/stress.go:25:15 cmd/seed.go:21:15
cmd/listing.go:73:52:1 other occurrence(s) of "job-start" found in: cmd/listing.go:92:52
cmd/listing.go:92:52:1 other occurrence(s) of "job-start" found in: cmd/listing.go:73:52
cmd/seed.go:46:24:1 other occurrence(s) of "DATABASE_URL" found in: cmd/listing.go:111:24
cmd/listing.go:111:24:1 other occurrence(s) of "DATABASE_URL" found in: cmd/seed.go:46:24
cmd/listing.go:66:52:1 other occurrence(s) of "website" found in: cmd/listing.go:85:52
cmd/listing.go:85:52:1 other occurrence(s) of "website" found in: cmd/listing.go:66:52
cmd/server.go:109:26:1 other occurrence(s) of "dev-secret-key" found in: cmd/server.go:112:33
cmd/server.go:112:33:1 other occurrence(s) of "dev-secret-key" found in: cmd/server.go:109:26
cmd/listing.go:76:52:1 other occurrence(s) of "pay-range" found in: cmd/listing.go:95:52
cmd/listing.go:95:52:1 other occurrence(s) of "pay-range" found in: cmd/listing.go:76:52
cmd/seed.go:39:12:2 other occurrence(s) of ".tester/data/agbalumo.db" found in: cmd/seed.go:45:15 cmd/listing.go:114:9
cmd/seed.go:45:15:2 other occurrence(s) of ".tester/data/agbalumo.db" found in: cmd/seed.go:39:12 cmd/listing.go:114:9
cmd/listing.go:114:9:2 other occurrence(s) of ".tester/data/agbalumo.db" found in: cmd/seed.go:39:12 cmd/seed.go:45:15
cmd/listing.go:70:54:1 other occurrence(s) of "event-start" found in: cmd/listing.go:89:54
cmd/listing.go:89:54:1 other occurrence(s) of "event-start" found in: cmd/listing.go:70:54
cmd/server.go:51:5:1 other occurrence(s) of "status" found in: cmd/server.go:193:50
cmd/server.go:193:50:1 other occurrence(s) of "status" found in: cmd/server.go:51:5
cmd/listing.go:63:50:1 other occurrence(s) of "email" found in: cmd/listing.go:82:50
cmd/listing.go:82:50:1 other occurrence(s) of "email" found in: cmd/listing.go:63:50
goconst: 2026/04/01 11:47:31 Found 0 Go files to process in batches of 50
```
## 3. Struct Alignment (fieldalignment)
```
-: no Go files in /Users/johnnyblase/gym/agbalumo/internal
fieldalignment: analysis skipped due to errors in package
/Users/johnnyblase/gym/agbalumo/cmd/serve.go:13:19: struct with 48 pointer bytes could be 40
/Users/johnnyblase/gym/agbalumo/cmd/seed_test.go:8:13: struct with 56 pointer bytes could be 48
/Users/johnnyblase/gym/agbalumo/cmd/serve_test.go:10:13: struct of size 136 could be 128
/Users/johnnyblase/gym/agbalumo/cmd/server_public_test.go:13:13: struct with 64 pointer bytes could be 56
```
## 4. Code Duplication
```
found 2 clones:
  internal/repository/sqlite/sqlite.go:38,46
  internal/repository/sqlite/sqlite.go:41,49
found 3 clones:
  internal/agent/verify_apispec_test.go:159,206
  internal/agent/verify_apispec_test.go:533,562
  internal/agent/verify_apispec_test.go:564,606
found 2 clones:
  internal/module/admin/admin.go:276,279
  internal/module/admin/admin_bulk.go:21,24
found 2 clones:
  internal/repository/cached/cached_test.go:77,94
  internal/repository/cached/cached_test.go:163,180
found 12 clones:
  internal/handler/feedback_integration_test.go:39,39
  internal/handler/feedback_integration_test.go:85,85
  internal/handler/feedback_integration_test.go:105,105
  internal/handler/feedback_integration_test.go:124,124
  internal/handler/feedback_integration_test.go:149,149
  internal/module/admin/admin_bulk_integration_test.go:49,49
  internal/module/admin/admin_bulk_integration_test.go:82,82
  internal/module/admin/admin_bulk_test.go:157,157
  internal/module/admin/admin_bulk_test.go:179,179
  internal/module/admin/admin_bulk_test.go:199,199
  internal/module/admin/admin_bulk_test.go:225,225
  internal/module/admin/admin_bulk_test.go:299,299
found 3 clones:
  internal/agent/verify_security_test.go:18,18
  internal/agent/verify_security_test.go:34,34
  internal/agent/verify_security_test.go:51,51
found 2 clones:
  cmd/listing_create.go:58,72
  cmd/listing_update.go:56,70
found 4 clones:
  cmd/harness/commands/chaos.go:37,39
  cmd/harness/commands/root.go:103,105
  internal/agent/progress.go:147,149
  internal/agent/progress.go:154,156
found 2 clones:
  internal/domain/job_validation_test.go:76,78
  internal/middleware/ratelimit_test.go:140,142
found 3 clones:
  internal/repository/sqlite/sqlite_category_test.go:222,224
  internal/repository/sqlite/sqlite_category_test.go:223,225
  internal/repository/sqlite/sqlite_category_test.go:224,226
found 2 clones:
  internal/ui/renderer.go:65,66
  internal/ui/renderer.go:66,67
found 17 clones:
  internal/repository/sqlite/feedback_test.go:87,89
  internal/repository/sqlite/feedback_test.go:174,176
  internal/repository/sqlite/feedback_test.go:177,179
  internal/repository/sqlite/feedback_test.go:180,182
  internal/repository/sqlite/repro_test.go:72,74
  internal/repository/sqlite/repro_test.go:78,80
  internal/repository/sqlite/repro_test.go:113,115
  internal/repository/sqlite/sqlite_category_test.go:28,30
  internal/repository/sqlite/sqlite_category_test.go:62,64
  internal/repository/sqlite/sqlite_category_test.go:70,72
  internal/repository/sqlite/sqlite_category_test.go:151,153
  internal/repository/sqlite/sqlite_category_test.go:229,231
  internal/repository/sqlite/sqlite_category_test.go:285,287
  internal/repository/sqlite/sqlite_event_test.go:29,31
  internal/repository/sqlite/sqlite_job_test.go:34,36
  internal/repository/sqlite/sqlite_listing_test.go:31,33
  internal/util/fs_test.go:168,170
found 2 clones:
  internal/module/admin/admin_dashboard_test.go:55,55
  internal/module/admin/admin_dashboard_test.go:56,56
found 3 clones:
  internal/module/listing/listing_form.go:65,65
  internal/module/listing/listing_mutations.go:136,136
  internal/module/listing/listing_mutations.go:181,181
found 2 clones:
  internal/module/listing/listing_form.go:95,101
  internal/module/listing/listing_form.go:102,108
found 3 clones:
  cmd/security-audit/audit_test.go:64,66
  cmd/security-audit/audit_test.go:183,185
  cmd/seed_test.go:52,54
found 8 clones:
  internal/handler/mock_repository_test.go:28,30
  internal/handler/mock_repository_test.go:68,70
  internal/handler/mock_repository_test.go:72,74
  internal/handler/mock_repository_test.go:108,110
  internal/module/admin/mock_repository_test.go:28,30
  internal/module/admin/mock_repository_test.go:68,70
  internal/module/admin/mock_repository_test.go:72,74
  internal/module/admin/mock_repository_test.go:108,110
found 5 clones:
  internal/module/admin/admin_bulk_test.go:44,44
  internal/module/admin/admin_bulk_test.go:79,79
  internal/module/admin/admin_bulk_test.go:94,94
  internal/module/admin/admin_bulk_test.go:123,123
  internal/module/admin/admin_bulk_test.go:269,269
found 2 clones:
  internal/agent/verify_apispec_test.go:273,311
  internal/agent/verify_apispec_test.go:355,399
found 22 clones:
  internal/module/auth/handler.go:200,202
  internal/module/auth/handler.go:209,211
  internal/module/auth/handler.go:217,219
  internal/module/auth/handler.go:243,245
  internal/module/auth/handler.go:256,258
  internal/module/auth/handler.go:261,263
  internal/module/auth/handler.go:266,268
  internal/module/auth/handler.go:305,308
  internal/module/auth/handler.go:321,323
  internal/module/listing/listing.go:208,210
  internal/module/listing/listing.go:227,229
  internal/module/listing/listing.go:232,234
  internal/module/listing/listing.go:237,239
  internal/module/listing/listing_mutations.go:33,35
  internal/module/listing/listing_mutations.go:41,43
  internal/module/listing/listing_mutations.go:58,60
  internal/module/listing/listing_mutations.go:65,67
  internal/module/listing/listing_mutations.go:70,72
  internal/module/listing/listing_mutations.go:103,105
  internal/module/listing/listing_mutations.go:115,117
  internal/module/listing/listing_mutations.go:121,123
  internal/module/listing/listing_mutations.go:125,127
found 2 clones:
  internal/module/admin/admin_dashboard_test.go:37,37
  internal/repository/sqlite/sqlite_listing_test.go:166,166
found 8 clones:
  internal/module/auth/handler_login_test.go:105,105
  internal/module/auth/handler_login_test.go:271,271
  internal/module/auth/handler_login_test.go:298,298
  internal/module/auth/handler_register_test.go:40,40
  internal/module/auth/handler_register_test.go:81,81
  internal/module/auth/handler_register_test.go:127,127
  internal/module/auth/handler_register_test.go:170,170
  internal/module/auth/handler_register_test.go:211,211
found 2 clones:
  internal/module/admin/admin_actions_test.go:63,74
  internal/module/admin/admin_actions_test.go:86,97
found 2 clones:
  internal/middleware/ratelimit.go:23,26
  internal/repository/sqlite/sqlite.go:15,18
found 8 clones:
  cmd/harness/commands/chaos.go:46,46
  internal/agent/ast.go:32,32
  internal/agent/cost.go:59,59
  internal/agent/drift.go:140,140
  internal/agent/drift.go:191,191
  internal/agent/security.go:97,97
  internal/agent/verify.go:434,434
  internal/module/listing/ui_regression_home_test.go:70,70
found 4 clones:
  internal/repository/sqlite/search_performance_test.go:64,66
  internal/repository/sqlite/search_performance_test.go:70,72
  internal/repository/sqlite/search_performance_test.go:88,90
  internal/repository/sqlite/sqlite_listing_test.go:474,476
found 4 clones:
  internal/handler/mock_repository_test.go:88,90
  internal/handler/mock_repository_test.go:132,134
  internal/module/admin/mock_repository_test.go:88,90
  internal/module/admin/mock_repository_test.go:132,134
found 4 clones:
  internal/module/auth/handler_login_test.go:65,65
  internal/module/auth/handler_login_test.go:94,94
  internal/module/auth/handler_login_test.go:139,139
  internal/module/auth/handler_register_test.go:201,201
found 2 clones:
  internal/module/auth/internal_test.go:39,48
  internal/module/auth/internal_test.go:50,59
found 4 clones:
  internal/handler/mock_repository_test.go:72,78
  internal/handler/mock_repository_test.go:108,114
  internal/module/admin/mock_repository_test.go:72,78
  internal/module/admin/mock_repository_test.go:108,114
found 3 clones:
  internal/handler/real_template_helpers_test.go:44,48
  internal/module/admin/ui_helpers_test.go:44,48
  internal/module/listing/ui_helpers_test.go:44,48
found 2 clones:
  internal/module/listing/listing.go:134,137
  internal/module/listing/pagination.go:66,69
found 4 clones:
  internal/domain/listing_validation_test.go:416,422
  internal/domain/listing_validation_test.go:423,429
  internal/domain/listing_validation_test.go:430,436
  internal/domain/listing_validation_test.go:437,443
found 2 clones:
  internal/module/listing/listing_edit_reproduction_test.go:80,82
  internal/module/listing/listing_edit_reproduction_test.go:149,151
found 3 clones:
  internal/common/page_handler.go:34,34
  internal/module/admin/admin.go:241,241
  internal/module/listing/listing.go:267,267
found 4 clones:
  cmd/harness/commands/root.go:17,22
  internal/agent/cost.go:17,22
  internal/agent/progress.go:10,15
  internal/domain/csv.go:9,14
found 6 clones:
  internal/common/page_handler.go:35,38
  internal/module/admin/admin.go:205,208
  internal/module/admin/admin.go:211,214
  internal/module/admin/admin_listings.go:41,44
  internal/module/listing/listing.go:124,127
  internal/module/listing/listing.go:268,271
found 3 clones:
  cmd/benchmark.go:37,45
  cmd/serve_test.go:10,18
  internal/service/geocoding_test.go:13,21
found 4 clones:
  internal/domain/url.go:11,11
  internal/handler/ui_regression_admin_test.go:155,155
  internal/handler/ui_regression_admin_test.go:159,159
  internal/module/listing/ui_regression_typography_test.go:23,23
found 4 clones:
  internal/service/csv.go:60,64
  internal/service/csv.go:68,72
  internal/service/csv.go:76,80
  internal/service/csv.go:129,133
found 3 clones:
  internal/repository/sqlite/sqlite_featured_test.go:98,106
  internal/repository/sqlite/sqlite_featured_test.go:101,109
  internal/repository/sqlite/sqlite_featured_test.go:104,112
found 9 clones:
  cmd/harness/commands/gate_test.go:14,17
  cmd/harness/commands/init_test.go:13,16
  cmd/harness/commands/root_test.go:11,14
  cmd/harness/commands/root_test.go:45,48
  cmd/harness/commands/root_test.go:98,101
  cmd/harness/commands/set_phase_test.go:14,17
  cmd/harness/commands/status_test.go:14,17
  cmd/harness/commands/update_coverage_test.go:14,17
  cmd/harness/commands/verify_test.go:16,19
found 2 clones:
  internal/agent/progress.go:45,48
  internal/agent/progress.go:72,75
found 3 clones:
  cmd/harness/commands/gate.go:28,28
  cmd/harness/commands/init.go:21,21
  cmd/harness/commands/verify.go:30,30
found 2 clones:
  cmd/security-audit/audit_test.go:55,60
  cmd/security-audit/audit_test.go:313,318
found 3 clones:
  internal/handler/mock_repository_test.go:40,40
  internal/module/admin/mock_repository_test.go:40,40
  internal/repository/sqlite/sqlite_listing_read.go:176,176
found 3 clones:
  internal/handler/mock_repository_test.go:104,104
  internal/module/admin/mock_repository_test.go:104,104
  internal/repository/sqlite/sqlite_category.go:30,30
found 2 clones:
  internal/service/csv_test.go:15,36
  internal/service/csv_test.go:51,78
found 2 clones:
  cmd/harness/commands/root_test.go:22,24
  internal/agent/cost_test.go:88,90
found 2 clones:
  cmd/harness/commands/chaos_test.go:57,59
  internal/service/csv_test.go:46,48
found 3 clones:
  internal/agent/verify_test.go:232,253
  internal/agent/verify_test.go:255,276
  internal/agent/verify_test.go:278,298
found 4 clones:
  internal/module/listing/listing_form.go:83,89
  internal/module/listing/listing_form.go:95,101
  internal/module/listing/listing_form.go:102,108
  internal/module/listing/listing_form.go:114,120
found 2 clones:
  internal/agent/security.go:352,352
  internal/agent/security.go:498,498
found 2 clones:
  cmd/harness/commands/init.go:45,47
  cmd/harness/commands/verify.go:156,158
found 2 clones:
  internal/module/admin/admin_all_listings_test.go:88,89
  internal/module/admin/admin_all_listings_test.go:89,90
found 5 clones:
  internal/module/listing/listing_read_test.go:41,51
  internal/module/listing/listing_read_test.go:81,91
  internal/module/listing/listing_read_test.go:122,132
  internal/module/listing/listing_read_test.go:163,173
  internal/module/listing/ui_regression_home_test.go:44,54
found 3 clones:
  internal/handler/real_template_helpers_test.go:30,30
  internal/module/admin/ui_helpers_test.go:30,30
  internal/module/listing/ui_helpers_test.go:30,30
found 2 clones:
  internal/module/admin/admin_bulk_integration_test.go:121,121
  internal/service/csv.go:25,25
found 4 clones:
  internal/module/admin/admin_all_listings_test.go:21,21
  internal/module/admin/admin_all_listings_test.go:88,88
  internal/module/admin/admin_all_listings_test.go:89,89
  internal/module/admin/admin_all_listings_test.go:90,90
found 2 clones:
  cmd/benchmark.go:54,54
  cmd/benchmark.go:59,59
found 2 clones:
  internal/domain/repository.go:23,23
  internal/domain/repository.go:67,67
found 2 clones:
  internal/module/listing/listing_service.go:41,43
  internal/module/listing/listing_service.go:50,52
found 15 clones:
  internal/module/listing/listing_featured_test.go:31,31
  internal/module/listing/listing_featured_test.go:32,32
  internal/module/listing/listing_featured_test.go:33,33
  internal/module/listing/listing_featured_test.go:34,34
  internal/module/listing/listing_featured_test.go:86,86
  internal/module/listing/listing_featured_test.go:87,87
  internal/module/listing/listing_featured_test.go:88,88
  internal/module/listing/listing_featured_test.go:89,89
  internal/module/listing/listing_featured_test.go:135,135
  internal/module/listing/listing_featured_test.go:136,136
  internal/module/listing/listing_featured_test.go:174,174
  internal/module/listing/listing_featured_test.go:175,175
  internal/module/listing/listing_featured_test.go:212,212
  internal/module/listing/listing_featured_test.go:213,213
  internal/module/listing/listing_featured_test.go:214,214
found 2 clones:
  internal/agent/state.go:82,84
  internal/module/auth/handler.go:131,133
found 2 clones:
  internal/module/listing/listing_mutations.go:103,105
  internal/module/listing/listing_mutations.go:125,127
found 2 clones:
  internal/agent/verify.go:223,225
  internal/agent/verify.go:367,369
found 2 clones:
  internal/repository/sqlite/feedback_test.go:174,179
  internal/repository/sqlite/feedback_test.go:177,182
found 2 clones:
  internal/handler/mock_geocoding_test.go:1,14
  internal/module/listing/mock_geocoding_test.go:1,14
found 4 clones:
  internal/repository/sqlite/sqlite.go:38,40
  internal/repository/sqlite/sqlite.go:41,43
  internal/repository/sqlite/sqlite.go:44,46
  internal/repository/sqlite/sqlite.go:47,49
found 2 clones:
  internal/module/listing/listing.go:100,103
  internal/module/listing/listing.go:108,111
found 2 clones:
  internal/service/image.go:190,203
  internal/service/image.go:206,219
found 8 clones:
  internal/domain/repository.go:13,13
  internal/handler/mock_repository_test.go:40,40
  internal/module/admin/mock_repository_test.go:40,40
  internal/module/auth/handler.go:37,37
  internal/module/auth/handler.go:100,100
  internal/module/auth/handler.go:152,152
  internal/module/auth/test_helpers_test.go:39,39
  internal/repository/sqlite/sqlite_listing_read.go:176,176
found 2 clones:
  cmd/seed_test.go:32,37
  cmd/seed_test.go:38,43
found 5 clones:
  internal/domain/image.go:10,10
  internal/handler/ui_helpers_test.go:47,47
  internal/module/listing/listing_helpers_test.go:69,69
  internal/seeder/category_seeder.go:16,16
  internal/service/image.go:48,48
found 5 clones:
  internal/domain/job_validation_test.go:39,45
  internal/domain/job_validation_test.go:46,52
  internal/domain/job_validation_test.go:53,59
  internal/domain/job_validation_test.go:60,66
  internal/domain/job_validation_test.go:89,95
found 3 clones:
  internal/module/listing/listing.go:237,237
  internal/module/listing/listing_mutations.go:70,70
  internal/seeder/seeder.go:45,45
found 23 clones:
  internal/module/admin/admin_actions_test.go:26,26
  internal/module/admin/admin_actions_test.go:118,118
  internal/module/admin/admin_actions_test.go:148,148
  internal/module/admin/admin_actions_test.go:165,165
  internal/module/admin/admin_all_listings_test.go:75,75
  internal/module/admin/admin_all_listings_test.go:95,95
  internal/module/admin/admin_bulk_test.go:244,244
  internal/module/admin/admin_category_test.go:32,32
  internal/module/admin/admin_category_test.go:52,52
  internal/module/admin/admin_category_test.go:76,76
  internal/module/admin/admin_category_test.go:99,99
  internal/module/admin/admin_claims_test.go:18,18
  internal/module/admin/admin_claims_test.go:39,39
  internal/module/admin/admin_claims_test.go:53,53
  internal/module/admin/admin_claims_test.go:74,74
  internal/module/admin/admin_dashboard_test.go:39,39
  internal/module/admin/admin_dashboard_test.go:58,58
  internal/module/admin/admin_featured_mock_test.go:29,29
  internal/module/admin/admin_featured_mock_test.go:45,45
  internal/module/admin/admin_login_test.go:51,51
  internal/module/admin/admin_middleware_test.go:16,16
  internal/module/admin/admin_middleware_test.go:32,32
  internal/module/admin/admin_users_test.go:26,26
found 3 clones:
  internal/service/image.go:148,148
  internal/service/image.go:197,197
  internal/service/image.go:213,213
found 3 clones:
  internal/module/listing/listing.go:172,172
  internal/module/listing/listing_claim.go:26,26
  internal/module/listing/listing_claim.go:28,28
found 2 clones:
  cmd/harness/commands/verify.go:100,104
  cmd/harness/commands/verify.go:107,111
found 10 clones:
  internal/module/auth/handler_login_test.go:54,54
  internal/module/auth/handler_login_test.go:84,84
  internal/module/auth/handler_login_test.go:235,235
  internal/module/auth/handler_login_test.go:261,261
  internal/module/auth/handler_login_test.go:287,287
  internal/module/auth/handler_register_test.go:28,28
  internal/module/auth/handler_register_test.go:55,55
  internal/module/auth/handler_register_test.go:101,101
  internal/module/auth/handler_register_test.go:144,144
  internal/module/auth/handler_register_test.go:187,187
found 2 clones:
  internal/agent/task_optimization_test.go:69,71
  internal/agent/task_optimization_test.go:72,74
found 5 clones:
  cmd/server.go:234,235
  internal/agent/security.go:110,110
  internal/agent/security.go:249,251
  internal/common/error_ui.go:18,20
  internal/module/listing/ui_regression_home_test.go:84,84
found 4 clones:
  internal/domain/listing.go:237,239
  internal/domain/listing.go:240,242
  internal/domain/listing.go:243,245
  internal/domain/listing.go:246,248
found 2 clones:
  internal/module/listing/pagination.go:41,43
  internal/ui/renderer.go:70,72
found 2 clones:
  internal/agent/verify_apispec_test.go:528,562
  internal/agent/verify_apispec_test.go:559,606
found 3 clones:
  internal/agent/state.go:125,127
  internal/util/fs.go:25,27
  internal/util/fs.go:47,49
found 2 clones:
  internal/agent/verify.go:465,474
  internal/agent/verify.go:483,491
found 16 clones:
  cmd/listing_backfill.go:31,31
  internal/module/listing/ui_regression_home_test.go:38,38
  internal/repository/sqlite/search_performance_test.go:65,65
  internal/repository/sqlite/search_performance_test.go:71,71
  internal/repository/sqlite/search_performance_test.go:89,89
  internal/repository/sqlite/sqlite_listing_bulk_test.go:29,29
  internal/repository/sqlite/sqlite_listing_test.go:70,70
  internal/repository/sqlite/sqlite_listing_test.go:79,79
  internal/repository/sqlite/sqlite_listing_test.go:85,85
  internal/repository/sqlite/sqlite_listing_test.go:91,91
  internal/repository/sqlite/sqlite_listing_test.go:97,97
  internal/repository/sqlite/sqlite_listing_test.go:113,113
  internal/repository/sqlite/sqlite_listing_test.go:127,127
  internal/repository/sqlite/sqlite_listing_test.go:475,475
  internal/repository/sqlite/sqlite_test.go:50,50
  internal/seeder/seeder.go:27,27
found 2 clones:
  internal/agent/coverage_test.go:86,98
  internal/agent/coverage_test.go:100,110
found 2 clones:
  internal/agent/verify_test.go:343,363
  internal/agent/verify_test.go:354,374
found 2 clones:
  internal/repository/sqlite/sqlite_listing_read.go:121,127
  internal/repository/sqlite/sqlite_listing_read.go:196,202
found 2 clones:
  internal/module/admin/admin_actions_test.go:48,54
  internal/module/admin/admin_actions_test.go:55,62
found 16 clones:
  internal/module/admin/admin_bulk_integration_test.go:24,24
  internal/module/admin/admin_bulk_integration_test.go:74,74
  internal/module/admin/admin_bulk_integration_test.go:99,99
  internal/module/admin/admin_bulk_integration_test.go:139,139
  internal/module/admin/admin_delete_test.go:34,34
  internal/module/admin/admin_delete_test.go:97,97
  internal/module/admin/admin_delete_test.go:117,117
  internal/module/admin/admin_delete_test.go:137,137
  internal/module/admin/admin_export_test.go:23,23
  internal/module/admin/admin_login_test.go:117,117
  internal/module/admin/admin_ui_integration_test.go:41,41
  internal/module/admin/admin_ui_integration_test.go:81,81
  internal/module/admin/admin_ui_integration_test.go:119,119
  internal/module/admin/admin_ui_integration_test.go:157,157
  internal/module/admin/admin_ui_integration_test.go:183,183
  internal/module/admin/admin_ui_integration_test.go:219,219
found 23 clones:
  cmd/stress_test.go:45,47
  internal/agent/verify_apispec_test.go:620,622
  internal/repository/sqlite/feedback_test.go:150,152
  internal/repository/sqlite/repro_test.go:33,35
  internal/repository/sqlite/sqlite_category_test.go:107,109
  internal/repository/sqlite/sqlite_category_test.go:116,118
  internal/repository/sqlite/sqlite_category_test.go:252,254
  internal/repository/sqlite/sqlite_listing_test.go:74,76
  internal/repository/sqlite/sqlite_listing_test.go:86,88
  internal/repository/sqlite/sqlite_listing_test.go:92,94
  internal/repository/sqlite/sqlite_listing_test.go:98,100
  internal/repository/sqlite/sqlite_listing_test.go:118,120
  internal/repository/sqlite/sqlite_listing_test.go:132,134
  internal/repository/sqlite/sqlite_listing_test.go:244,246
  internal/repository/sqlite/sqlite_listing_test.go:273,275
  internal/repository/sqlite/sqlite_listing_test.go:359,361
  internal/repository/sqlite/sqlite_listing_test.go:387,389
  internal/repository/sqlite/sqlite_listing_test.go:393,395
  internal/repository/sqlite/sqlite_user_test.go:40,42
  internal/repository/sqlite/sqlite_user_test.go:62,64
  internal/seeder/category_seeder_test.go:51,53
  internal/seeder/category_seeder_test.go:102,104
  internal/seeder/config_verification_test.go:62,64
found 3 clones:
  internal/handler/mock_repository_test.go:92,92
  internal/module/admin/mock_repository_test.go:92,92
  internal/repository/sqlite/sqlite_user.go:89,89
found 5 clones:
  internal/agent/state_test.go:46,51
  internal/service/csv_test.go:27,32
  internal/service/csv_test.go:30,35
  internal/service/csv_test.go:69,74
  internal/service/csv_test.go:72,77
found 2 clones:
  internal/agent/security.go:220,228
  internal/agent/security.go:591,599
found 2 clones:
  internal/domain/repository.go:14,14
  internal/domain/repository.go:54,54
found 2 clones:
  cmd/server.go:234,237
  internal/agent/security.go:249,256
found 2 clones:
  cmd/harness/commands/verify.go:46,46
  cmd/harness/commands/verify.go:142,142
found 5 clones:
  cmd/harness/commands/chaos.go:82,84
  cmd/harness/commands/gate.go:19,19
  internal/agent/security.go:116,116
  internal/agent/verify.go:466,466
  internal/agent/verify.go:484,484
found 2 clones:
  internal/repository/sqlite/sqlite_stats.go:53,61
  internal/repository/sqlite/sqlite_stats.go:64,72
found 2 clones:
  internal/agent/verify_test.go:186,203
  internal/agent/verify_test.go:205,222
found 9 clones:
  internal/agent/verify_apispec_test.go:242,251
  internal/agent/verify_apispec_test.go:273,282
  internal/agent/verify_apispec_test.go:301,310
  internal/agent/verify_apispec_test.go:329,338
  internal/agent/verify_apispec_test.go:355,364
  internal/agent/verify_apispec_test.go:389,398
  internal/agent/verify_apispec_test.go:427,436
  internal/agent/verify_apispec_test.go:470,479
  internal/agent/verify_apispec_test.go:516,526
found 2 clones:
  internal/repository/sqlite/sqlite_user.go:57,66
  internal/repository/sqlite/sqlite_user.go:69,78
found 4 clones:
  internal/handler/mock_repository_test.go:48,50
  internal/handler/mock_repository_test.go:84,86
  internal/module/admin/mock_repository_test.go:48,50
  internal/module/admin/mock_repository_test.go:84,86
found 5 clones:
  internal/handler/ui_regression_admin_test.go:198,200
  internal/handler/ui_regression_modals_test.go:66,68
  internal/module/admin/admin_ui_integration_test.go:235,237
  internal/module/listing/listing_card_test.go:42,44
  internal/module/listing/ui_regression_home_test.go:115,117
found 2 clones:
  internal/repository/sqlite/sqlite_listing_test.go:461,471
  internal/repository/sqlite/sqlite_listing_test.go:486,496
found 9 clones:
  internal/handler/mock_geocoding_test.go:9,9
  internal/handler/mock_repository_test.go:36,36
  internal/module/admin/mock_repository_test.go:36,36
  internal/module/listing/listing_geocoding_test.go:21,21
  internal/module/listing/mock_geocoding_test.go:9,9
  internal/repository/sqlite/sqlite_listing_read.go:207,207
  internal/service/csv_test.go:146,146
  internal/service/csv_test.go:177,177
  internal/service/geocoding.go:34,34
found 2 clones:
  internal/module/auth/middleware_test.go:85,91
  internal/module/auth/middleware_test.go:108,114
found 2 clones:
  internal/agent/verify_test.go:313,321
  internal/agent/verify_test.go:323,331
found 2 clones:
  internal/common/page_handler.go:32,44
  internal/module/listing/listing.go:265,277
found 2 clones:
  internal/repository/cached/cached.go:34,34
  internal/repository/cached/cached.go:69,69
found 2 clones:
  internal/repository/sqlite/sqlite_listing_read.go:82,82
  internal/repository/sqlite/sqlite_listing_read.go:210,210
found 4 clones:
  internal/repository/sqlite/sqlite_featured_test.go:98,103
  internal/repository/sqlite/sqlite_featured_test.go:101,106
  internal/repository/sqlite/sqlite_featured_test.go:104,109
  internal/repository/sqlite/sqlite_featured_test.go:107,112
found 2 clones:
  internal/agent/state_test.go:126,138
  internal/agent/state_test.go:140,152
found 4 clones:
  internal/module/admin/admin_delete_test.go:27,27
  internal/module/admin/admin_delete_test.go:49,49
  internal/module/admin/admin_delete_test.go:115,115
  internal/module/admin/admin_delete_test.go:129,129
found 9 clones:
  cmd/harness/commands/chaos.go:82,83
  cmd/harness/commands/gate.go:19,19
  cmd/harness/commands/set_phase.go:16,16
  internal/agent/security.go:116,116
  internal/agent/security.go:318,318
  internal/agent/security.go:410,410
  internal/agent/verify.go:466,466
  internal/agent/verify.go:484,484
  internal/service/csv.go:256,256
found 2 clones:
  internal/ui/renderer.go:135,139
  internal/ui/renderer.go:141,145
found 2 clones:
  internal/service/csv_test.go:92,94
  internal/service/csv_test.go:161,163
found 2 clones:
  internal/repository/sqlite/sqlite_category_test.go:18,30
  internal/repository/sqlite/sqlite_category_test.go:52,64
found 2 clones:
  internal/module/admin/admin_helpers_test.go:28,28
  internal/module/listing/listing_helpers_test.go:44,44
found 5 clones:
  internal/service/csv.go:88,96
  internal/service/csv.go:91,99
  internal/service/csv.go:94,102
  internal/service/csv.go:97,105
  internal/service/csv.go:100,108
found 2 clones:
  cmd/security-audit/audit_test.go:377,381
  internal/agent/archive_test.go:95,99
found 4 clones:
  internal/module/listing/listing_create_test.go:25,32
  internal/module/listing/listing_create_test.go:40,45
  internal/module/listing/listing_create_test.go:46,51
  internal/module/listing/listing_form_integration_test.go:68,73
found 14 clones:
  cmd/listing_cmd_test.go:25,25
  cmd/listing_cmd_test.go:57,57
  cmd/listing_cmd_test.go:87,87
  cmd/listing_cmd_test.go:102,102
  cmd/listing_cmd_test.go:166,166
  cmd/listing_read.go:25,25
  cmd/stress_test.go:40,40
  internal/module/admin/admin_bulk_test.go:128,128
  internal/module/listing/listing_form_integration_test.go:31,31
  internal/module/listing/listing_form_integration_test.go:46,46
  internal/module/listing/listing_form_integration_test.go:60,60
  internal/seeder/seeder_test.go:20,20
  internal/seeder/seeder_test.go:31,31
  internal/seeder/seeder_test.go:45,45
found 2 clones:
  internal/repository/sqlite/sqlite_featured_test.go:55,58
  internal/repository/sqlite/sqlite_featured_test.go:115,118
found 4 clones:
  internal/repository/sqlite/sqlite_category_test.go:222,223
  internal/repository/sqlite/sqlite_category_test.go:223,224
  internal/repository/sqlite/sqlite_category_test.go:224,225
  internal/repository/sqlite/sqlite_category_test.go:225,226
found 3 clones:
  internal/module/admin/admin_actions_test.go:68,70
  internal/module/admin/admin_actions_test.go:80,82
  internal/module/admin/admin_actions_test.go:91,93
found 3 clones:
  internal/repository/sqlite/repro_test.go:19,20
  internal/repository/sqlite/sqlite_category_test.go:266,267
  internal/repository/sqlite/sqlite_category_test.go:267,268
found 8 clones:
  internal/module/listing/listing_featured_test.go:31,32
  internal/module/listing/listing_featured_test.go:33,34
  internal/module/listing/listing_featured_test.go:86,87
  internal/module/listing/listing_featured_test.go:88,89
  internal/module/listing/listing_featured_test.go:135,136
  internal/module/listing/listing_featured_test.go:174,175
  internal/module/listing/listing_featured_test.go:212,213
  internal/module/listing/listing_featured_test.go:213,214
found 3 clones:
  internal/handler/ui_helpers_test.go:19,40
  internal/module/auth/test_helpers_test.go:19,27
  internal/module/listing/listing_helpers_test.go:21,42
found 2 clones:
  internal/agent/redtest_test.go:96,98
  internal/handler/response_test.go:34,36
found 2 clones:
  internal/service/image_test.go:266,266
  internal/service/image_test.go:374,374
found 2 clones:
  internal/agent/security.go:321,330
  internal/agent/security.go:545,554
found 2 clones:
  internal/agent/drift_test.go:49,52
  internal/agent/drift_test.go:65,68
found 3 clones:
  internal/service/geocoding_test.go:22,52
  internal/service/geocoding_test.go:59,81
  internal/service/geocoding_test.go:135,157
found 3 clones:
  internal/agent/verify.go:31,34
  internal/agent/verify.go:215,218
  internal/agent/verify.go:397,400
found 15 clones:
  internal/module/admin/admin_bulk_integration_test.go:113,113
  internal/module/admin/admin_bulk_test.go:54,54
  internal/module/admin/admin_bulk_test.go:143,143
  internal/module/admin/admin_category_test.go:36,36
  internal/module/admin/admin_category_test.go:56,56
  internal/module/admin/admin_category_test.go:80,80
  internal/module/admin/admin_delete_test.go:71,71
  internal/module/admin/admin_delete_test.go:101,101
  internal/module/auth/handler_login_test.go:114,114
  internal/module/auth/handler_login_test.go:147,147
  internal/module/auth/handler_login_test.go:171,171
  internal/module/auth/handler_logout_test.go:35,35
  internal/module/auth/handler_logout_test.go:54,54
  internal/module/auth/handler_register_test.go:218,218
  internal/module/auth/middleware_test.go:39,39
found 2 clones:
  internal/module/admin/admin_login_test.go:46,48
  internal/module/listing/listing_delete_test.go:71,73
found 2 clones:
  cmd/security-audit/audit_test.go:160,166
  internal/domain/listing_validation_test.go:409,415
found 2 clones:
  internal/repository/sqlite/sqlite_claim.go:42,42
  internal/repository/sqlite/sqlite_claim.go:98,98
found 2 clones:
  internal/module/listing/listing.go:148,148
  internal/module/listing/listing.go:187,187
found 7 clones:
  internal/repository/cached/cached_test.go:49,51
  internal/repository/cached/cached_test.go:52,54
  internal/repository/cached/cached_test.go:72,74
  internal/repository/cached/cached_test.go:121,123
  internal/repository/sqlite/feedback_test.go:189,191
  internal/repository/sqlite/feedback_test.go:192,194
  internal/repository/sqlite/feedback_test.go:195,197
found 7 clones:
  internal/handler/integration_ds_test.go:45,45
  internal/handler/integration_ds_test.go:105,105
  internal/module/listing/listing_event_test.go:42,42
  internal/module/listing/listing_geocoding_test.go:45,45
  internal/module/listing/listing_update_image_test.go:60,60
  internal/module/listing/listing_upload_test.go:60,60
  internal/module/listing/listing_upload_test.go:107,107
found 2 clones:
  internal/agent/verify_apispec_test.go:29,41
  internal/agent/verify_apispec_test.go:34,46
found 4 clones:
  internal/seeder/seeder.go:59,82
  internal/seeder/seeder.go:86,109
  internal/seeder/seeder.go:99,122
  internal/seeder/seeder.go:112,135
found 2 clones:
  internal/handler/integration_ds_test.go:26,26
  internal/module/listing/edit_job_test.go:42,42
found 2 clones:
  internal/service/image.go:100,100
  internal/service/image.go:112,112
found 2 clones:
  internal/module/admin/admin_claims_test.go:16,49
  internal/module/admin/admin_claims_test.go:51,83
found 2 clones:
  cmd/listing.go:158,158
  internal/agent/verify.go:42,42
found 3 clones:
  internal/module/admin/admin_bulk_integration_test.go:121,121
  internal/module/listing/listing_service.go:20,24
  internal/service/csv.go:25,25
found 5 clones:
  internal/repository/sqlite/sqlite_category_test.go:220,220
  internal/repository/sqlite/sqlite_listing_test.go:145,145
  internal/repository/sqlite/sqlite_listing_test.go:146,146
  internal/repository/sqlite/sqlite_listing_test.go:147,147
  internal/repository/sqlite/sqlite_listing_test.go:148,148
found 2 clones:
  internal/repository/sqlite/sqlite_category_test.go:205,205
  internal/repository/sqlite/sqlite_category_test.go:210,210
found 3 clones:
  cmd/listing_create.go:63,72
  cmd/listing_update.go:56,65
  cmd/listing_update.go:61,70
found 2 clones:
  internal/ui/renderer_test.go:259,259
  internal/ui/renderer_test.go:369,369
found 3 clones:
  internal/module/auth/handler.go:38,38
  internal/module/auth/handler.go:107,107
  internal/module/auth/handler.go:156,156
found 11 clones:
  internal/module/admin/admin_actions_test.go:68,68
  internal/module/admin/admin_actions_test.go:69,69
  internal/module/admin/admin_actions_test.go:70,70
  internal/module/admin/admin_actions_test.go:71,71
  internal/module/admin/admin_actions_test.go:80,80
  internal/module/admin/admin_actions_test.go:81,81
  internal/module/admin/admin_actions_test.go:82,82
  internal/module/admin/admin_actions_test.go:91,91
  internal/module/admin/admin_actions_test.go:92,92
  internal/module/admin/admin_actions_test.go:93,93
  internal/module/admin/admin_actions_test.go:94,94
found 2 clones:
  internal/module/listing/pagination_test.go:32,55
  internal/module/listing/pagination_test.go:40,63
found 4 clones:
  internal/module/auth/handler_register_test.go:72,78
  internal/module/auth/handler_register_test.go:118,124
  internal/module/auth/handler_register_test.go:161,167
  internal/seeder/seeder_test.go:39,39
found 2 clones:
  cmd/harness/commands/chaos_test.go:47,49
  internal/agent/state_test.go:276,278
found 3 clones:
  internal/agent/verify_test.go:42,44
  internal/agent/verify_test.go:96,98
  internal/agent/verify_test.go:508,508
found 5 clones:
  internal/module/listing/listing_featured_test.go:31,33
  internal/module/listing/listing_featured_test.go:32,34
  internal/module/listing/listing_featured_test.go:86,88
  internal/module/listing/listing_featured_test.go:87,89
  internal/module/listing/listing_featured_test.go:212,214
found 3 clones:
  internal/repository/sqlite/sqlite_listing_test.go:319,319
  internal/repository/sqlite/sqlite_listing_test.go:320,320
  internal/repository/sqlite/sqlite_user_test.go:56,56
found 2 clones:
  internal/agent/drift.go:196,203
  internal/agent/drift.go:212,219
found 2 clones:
  internal/module/admin/admin_ui_integration_test.go:28,39
  internal/repository/sqlite/sqlite_user_test.go:16,17
found 2 clones:
  internal/repository/sqlite/sqlite_listing_test.go:191,191
  internal/repository/sqlite/sqlite_listing_test.go:193,193
found 2 clones:
  internal/agent/ast.go:170,175
  internal/agent/drift.go:115,120
found 5 clones:
  internal/agent/archive_test.go:50,52
  internal/agent/archive_test.go:60,62
  internal/agent/redtest_test.go:28,30
  internal/agent/redtest_test.go:53,55
  internal/agent/redtest_test.go:88,90
found 4 clones:
  internal/module/admin/admin_category_test.go:39,39
  internal/module/admin/admin_category_test.go:59,59
  internal/module/admin/admin_category_test.go:83,83
  internal/module/admin/admin_category_test.go:105,105
found 3 clones:
  internal/module/listing/listing_service_test.go:30,30
  internal/module/listing/listing_service_test.go:67,67
  internal/module/listing/listing_service_test.go:90,90
found 2 clones:
  internal/agent/state_test.go:185,187
  internal/agent/state_test.go:290,292
found 2 clones:
  internal/domain/listing_job_test.go:86,101
  internal/domain/listing_job_test.go:102,117
found 2 clones:
  cmd/security-audit/audit_test.go:27,42
  cmd/security-audit/audit_test.go:35,50
found 2 clones:
  internal/module/admin/admin_login_test.go:55,57
  internal/module/admin/admin_login_test.go:121,123
found 13 clones:
  internal/module/admin/admin_actions_test.go:52,52
  internal/module/admin/admin_actions_test.go:59,60
  internal/module/listing/listing_create_test.go:28,30
  internal/module/listing/listing_create_test.go:36,36
  internal/module/listing/listing_create_test.go:43,43
  internal/module/listing/listing_create_test.go:49,49
  internal/module/listing/listing_delete_test.go:34,34
  internal/module/listing/listing_delete_test.go:40,41
  internal/module/listing/listing_delete_test.go:55,57
  internal/module/listing/listing_form_integration_test.go:28,28
  internal/module/listing/listing_form_integration_test.go:43,43
  internal/module/listing/listing_form_integration_test.go:57,57
  internal/module/listing/listing_form_integration_test.go:71,71
found 2 clones:
  internal/seeder/category_seeder_test.go:73,75
  internal/ui/renderer_test.go:23,25
found 8 clones:
  internal/module/listing/listing_delete_test.go:27,27
  internal/module/listing/listing_delete_test.go:48,48
  internal/module/listing/listing_update_test.go:28,28
  internal/module/listing/listing_update_test.go:36,36
  internal/module/listing/listing_update_test.go:81,81
  internal/module/listing/listing_update_test.go:90,90
  internal/module/listing/listing_update_test.go:166,166
  internal/module/listing/listing_update_test.go:167,167
found 2 clones:
  internal/repository/cached/cached_test.go:135,140
  internal/repository/cached/cached_test.go:155,160
found 2 clones:
  internal/agent/redtest.go:96,101
  internal/agent/redtest.go:103,108
found 5 clones:
  internal/module/admin/admin_actions_test.go:146,146
  internal/module/admin/admin_claims_test.go:21,21
  internal/module/admin/admin_claims_test.go:56,56
  internal/module/admin/admin_dashboard_test.go:31,31
  internal/module/admin/admin_users_test.go:22,22
found 12 clones:
  internal/agent/verify_coverage_test.go:16,18
  internal/agent/verify_coverage_test.go:55,57
  internal/module/admin/admin_ui_integration_test.go:47,49
  internal/module/admin/admin_ui_integration_test.go:87,89
  internal/module/admin/admin_ui_integration_test.go:125,127
  internal/module/admin/admin_ui_integration_test.go:170,172
  internal/module/admin/admin_ui_integration_test.go:225,227
  internal/module/listing/listing_featured_test.go:52,54
  internal/module/listing/listing_featured_test.go:107,109
  internal/module/listing/listing_featured_test.go:152,154
  internal/module/listing/listing_featured_test.go:191,193
  internal/module/listing/listing_featured_test.go:231,233
found 6 clones:
  internal/common/page_handler_test.go:42,44
  internal/module/listing/listing_edge_cases_test.go:58,60
  internal/module/listing/listing_read_test.go:49,51
  internal/module/listing/listing_read_test.go:89,91
  internal/module/listing/listing_read_test.go:130,132
  internal/module/listing/listing_read_test.go:171,173
found 3 clones:
  internal/repository/sqlite/search_performance_test.go:63,67
  internal/repository/sqlite/search_performance_test.go:69,73
  internal/repository/sqlite/search_performance_test.go:87,91
found 2 clones:
  internal/repository/sqlite/sqlite_featured_test.go:98,109
  internal/repository/sqlite/sqlite_featured_test.go:101,112
found 33 clones:
  internal/handler/integration_ds_test.go:32,39
  internal/handler/integration_ds_test.go:92,99
  internal/module/listing/edit_job_test.go:46,53
  internal/module/listing/listing_claim_test.go:29,36
  internal/module/listing/listing_create_test.go:62,69
  internal/module/listing/listing_create_test.go:88,95
  internal/module/listing/listing_create_test.go:109,116
  internal/module/listing/listing_delete_test.go:80,87
  internal/module/listing/listing_edge_cases_test.go:88,95
  internal/module/listing/listing_edge_cases_test.go:120,127
  internal/module/listing/listing_edge_cases_test.go:138,145
  internal/module/listing/listing_edit_reproduction_test.go:39,46
  internal/module/listing/listing_edit_reproduction_test.go:109,116
  internal/module/listing/listing_event_test.go:22,29
  internal/module/listing/listing_featured_test.go:43,50
  internal/module/listing/listing_featured_test.go:98,105
  internal/module/listing/listing_featured_test.go:143,150
  internal/module/listing/listing_featured_test.go:182,189
  internal/module/listing/listing_featured_test.go:222,229
  internal/module/listing/listing_form_integration_test.go:83,90
  internal/module/listing/listing_form_integration_test.go:112,119
  internal/module/listing/listing_read_test.go:41,48
  internal/module/listing/listing_read_test.go:81,88
  internal/module/listing/listing_read_test.go:122,129
  internal/module/listing/listing_read_test.go:163,170
  internal/module/listing/listing_update_image_test.go:45,52
  internal/module/listing/listing_update_test.go:53,60
  internal/module/listing/listing_update_test.go:109,116
  internal/module/listing/listing_update_test.go:133,140
  internal/module/listing/listing_update_test.go:153,160
  internal/module/listing/listing_update_test.go:178,185
  internal/module/listing/listing_upload_test.go:25,32
  internal/module/listing/ui_regression_home_test.go:44,51
found 4 clones:
  internal/service/csv.go:88,99
  internal/service/csv.go:91,102
  internal/service/csv.go:94,105
  internal/service/csv.go:97,108
found 3 clones:
  internal/repository/sqlite/sqlite_listing_read.go:166,172
  internal/repository/sqlite/sqlite_listing_read.go:285,291
  internal/repository/sqlite/sqlite_user.go:98,104
found 10 clones:
  internal/repository/sqlite/repro_test.go:29,29
  internal/repository/sqlite/repro_test.go:54,54
  internal/repository/sqlite/sqlite_category_test.go:103,103
  internal/repository/sqlite/sqlite_category_test.go:112,112
  internal/repository/sqlite/sqlite_category_test.go:248,248
  internal/repository/sqlite/sqlite_category_test.go:290,290
  internal/seeder/category_seeder_test.go:47,47
  internal/seeder/category_seeder_test.go:98,98
  internal/seeder/config_verification_test.go:34,34
  internal/seeder/config_verification_test.go:57,57
found 2 clones:
  internal/module/listing/cache_busting_integration_test.go:39,40
  internal/module/listing/listing_update_image_test.go:41,41
found 3 clones:
  internal/module/admin/admin_actions_test.go:24,26
  internal/module/admin/admin_actions_test.go:163,165
  internal/module/admin/admin_category_test.go:74,76
found 4 clones:
  internal/module/listing/pagination_test.go:32,39
  internal/module/listing/pagination_test.go:40,47
  internal/module/listing/pagination_test.go:48,55
  internal/module/listing/pagination_test.go:56,63
found 2 clones:
  internal/service/csv.go:164,167
  internal/service/csv.go:196,199
found 4 clones:
  internal/module/admin/admin_all_listings_test.go:73,75
  internal/module/admin/admin_all_listings_test.go:93,95
  internal/module/admin/admin_featured_mock_test.go:27,29
  internal/module/admin/admin_featured_mock_test.go:43,45
found 2 clones:
  internal/domain/listing_validation_test.go:265,269
  internal/domain/listing_validation_test.go:399,403
found 2 clones:
  internal/module/admin/admin_delete_test.go:32,34
  internal/module/admin/admin_delete_test.go:135,137
found 4 clones:
  cmd/security-audit/audit_test.go:132,132
  cmd/security-audit/audit_test.go:275,275
  cmd/security-audit/main.go:15,15
  cmd/security-audit/main.go:20,20
found 2 clones:
  internal/service/image_test.go:102,104
  internal/service/image_test.go:284,286
found 2 clones:
  cmd/harness/commands/chaos_test.go:97,99
  internal/handler/ui_regression_admin_test.go:183,185
found 2 clones:
  internal/module/auth/auth.go:9,13
  internal/module/auth/handler.go:167,171
found 2 clones:
  internal/module/listing/listing_form.go:82,91
  internal/module/listing/listing_form.go:113,122
found 3 clones:
  internal/repository/sqlite/sqlite_event_test.go:40,42
  internal/repository/sqlite/sqlite_event_test.go:43,45
  internal/repository/sqlite/sqlite_job_test.go:54,56
found 4 clones:
  internal/domain/image.go:10,10
  internal/handler/ui_helpers_test.go:47,47
  internal/module/listing/listing_helpers_test.go:69,69
  internal/service/image.go:48,48
found 5 clones:
  cmd/serve_test.go:19,34
  cmd/serve_test.go:27,42
  cmd/serve_test.go:35,50
  cmd/serve_test.go:43,58
  cmd/serve_test.go:51,66
found 2 clones:
  internal/module/admin/admin_bulk_test.go:216,217
  internal/module/admin/admin_bulk_test.go:289,290
found 3 clones:
  internal/agent/progress_test.go:23,28
  internal/agent/progress_test.go:29,34
  internal/agent/progress_test.go:51,56
found 2 clones:
  internal/domain/repository.go:11,12
  internal/domain/repository.go:79,80
found 6 clones:
  internal/module/listing/listing_delete_test.go:26,28
  internal/module/listing/listing_delete_test.go:47,49
  internal/module/listing/listing_update_test.go:27,29
  internal/module/listing/listing_update_test.go:35,37
  internal/module/listing/listing_update_test.go:80,82
  internal/module/listing/listing_update_test.go:89,91
found 2 clones:
  internal/agent/drift.go:30,33
  internal/agent/drift.go:55,58
found 2 clones:
  internal/module/admin/admin_login_test.go:76,76
  internal/module/admin/admin_login_test.go:91,91
found 2 clones:
  cmd/security-audit/main.go:187,192
  cmd/security-audit/main.go:207,212
found 2 clones:
  internal/agent/drift.go:127,132
  internal/util/slices.go:12,17
found 3 clones:
  internal/module/auth/handler_login_test.go:97,106
  internal/module/auth/handler_login_test.go:296,299
  internal/module/auth/handler_register_test.go:204,212
found 3 clones:
  internal/module/listing/listing_form.go:82,82
  internal/module/listing/listing_form.go:93,93
  internal/module/listing/listing_form.go:113,113
found 2 clones:
  internal/agent/ast.go:66,73
  internal/service/geocoding.go:98,105
found 3 clones:
  internal/repository/sqlite/sqlite_listing_test.go:235,235
  internal/repository/sqlite/sqlite_listing_test.go:236,236
  internal/repository/sqlite/sqlite_listing_test.go:237,237
found 2 clones:
  cmd/aglog/main.go:62,64
  internal/seeder/category_seeder.go:30,32
found 3 clones:
  internal/agent/drift_test.go:74,76
  internal/agent/drift_test.go:182,184
  internal/repository/sqlite/sqlite_category_test.go:239,241
found 2 clones:
  internal/agent/security.go:38,44
  internal/agent/security.go:48,54
found 4 clones:
  internal/agent/verify_test.go:343,352
  internal/agent/verify_test.go:354,363
  internal/agent/verify_test.go:365,374
  internal/agent/verify_test.go:387,396
found 2 clones:
  cmd/serve.go:13,18
  internal/agent/verify_apispec_test.go:53,58
found 2 clones:
  internal/agent/state.go:22,31
  internal/domain/brand.go:4,13
found 2 clones:
  internal/repository/sqlite/sqlite_category_test.go:222,225
  internal/repository/sqlite/sqlite_category_test.go:223,226
found 3 clones:
  internal/service/csv.go:60,64
  internal/service/csv.go:68,72
  internal/service/csv.go:76,80
found 2 clones:
  internal/module/listing/listing_delete_test.go:37,43
  internal/module/listing/listing_delete_test.go:52,59
found 2 clones:
  internal/repository/sqlite/feedback_test.go:170,171
  internal/repository/sqlite/feedback_test.go:171,172
found 3 clones:
  internal/module/admin/admin_actions_test.go:146,146
  internal/module/admin/admin_claims_test.go:21,21
  internal/module/admin/admin_claims_test.go:56,56
found 5 clones:
  internal/agent/verify_test.go:191,196
  internal/agent/verify_test.go:210,215
  internal/agent/verify_test.go:237,242
  internal/agent/verify_test.go:260,265
  internal/agent/verify_test.go:283,288
found 2 clones:
  cmd/server_admin_test.go:44,48
  cmd/server_user_test.go:54,58
found 2 clones:
  cmd/server.go:128,130
  internal/util/fs.go:41,44
found 3 clones:
  internal/domain/repository.go:11,11
  internal/domain/repository.go:16,16
  internal/domain/repository.go:79,79
found 2 clones:
  internal/module/listing/cache_busting_integration_test.go:29,36
  internal/module/listing/listing_geocoding_test.go:31,38
found 2 clones:
  internal/repository/sqlite/sqlite_listing_test.go:235,236
  internal/repository/sqlite/sqlite_listing_test.go:236,237
found 9 clones:
  internal/handler/ui_regression_admin_test.go:32,38
  internal/handler/ui_regression_admin_test.go:60,66
  internal/handler/ui_regression_admin_test.go:84,90
  internal/handler/ui_regression_admin_test.go:107,113
  internal/handler/ui_regression_modals_test.go:62,68
  internal/module/admin/admin_ui_integration_test.go:94,101
  internal/module/admin/admin_ui_integration_test.go:132,139
  internal/module/listing/listing_card_test.go:39,44
  internal/module/listing/ui_regression_home_test.go:111,117
found 7 clones:
  internal/service/csv.go:88,90
  internal/service/csv.go:91,93
  internal/service/csv.go:94,96
  internal/service/csv.go:97,99
  internal/service/csv.go:100,102
  internal/service/csv.go:103,105
  internal/service/csv.go:106,108
found 34 clones:
  cmd/cli_json_test.go:84,86
  cmd/harness/commands/root_test.go:38,40
  internal/agent/cost_test.go:70,72
  internal/agent/cost_test.go:84,86
  internal/agent/redtest_test.go:58,60
  internal/agent/redtest_test.go:92,94
  internal/agent/state_test.go:46,48
  internal/agent/state_test.go:49,51
  internal/agent/state_test.go:256,258
  internal/domain/category_test.go:21,23
  internal/repository/cached/cached_test.go:46,48
  internal/repository/cached/cached_test.go:69,71
  internal/repository/cached/cached_test.go:91,93
  internal/repository/cached/cached_test.go:135,137
  internal/repository/cached/cached_test.go:155,157
  internal/repository/cached/cached_test.go:177,179
  internal/repository/sqlite/sqlite_category_test.go:37,39
  internal/repository/sqlite/sqlite_category_test.go:78,80
  internal/repository/sqlite/sqlite_category_test.go:159,161
  internal/repository/sqlite/sqlite_category_test.go:181,183
  internal/repository/sqlite/sqlite_listing_test.go:40,42
  internal/repository/sqlite/sqlite_listing_test.go:55,57
  internal/repository/sqlite/sqlite_listing_test.go:309,311
  internal/repository/sqlite/sqlite_listing_test.go:567,569
  internal/repository/sqlite/sqlite_listing_test.go:575,577
  internal/repository/sqlite/sqlite_user_test.go:115,117
  internal/service/csv_test.go:27,29
  internal/service/csv_test.go:30,32
  internal/service/csv_test.go:33,35
  internal/service/csv_test.go:69,71
  internal/service/csv_test.go:72,74
  internal/service/csv_test.go:75,77
  internal/service/csv_test.go:95,97
  internal/service/csv_test.go:119,121
found 2 clones:
  internal/repository/sqlite/search_performance_test.go:75,79
  internal/repository/sqlite/search_performance_test.go:81,85
found 3 clones:
  cmd/security-audit/audit_test.go:27,34
  cmd/security-audit/audit_test.go:35,42
  cmd/security-audit/audit_test.go:43,50
found 2 clones:
  internal/agent/verify_apispec_test.go:407,421
  internal/agent/verify_apispec_test.go:444,455
found 2 clones:
  internal/module/admin/admin.go:317,317
  internal/module/admin/admin.go:318,318
found 8 clones:
  internal/middleware/ratelimit_test.go:22,24
  internal/middleware/ratelimit_test.go:42,44
  internal/middleware/ratelimit_test.go:73,75
  internal/middleware/ratelimit_test.go:92,94
  internal/middleware/ratelimit_test.go:126,128
  internal/middleware/security_test.go:19,21
  internal/module/auth/middleware_test.go:31,33
  internal/module/auth/middleware_test.go:56,58
found 9 clones:
  internal/handler/response_test.go:37,39
  internal/repository/sqlite/sqlite_job_test.go:45,47
  internal/repository/sqlite/sqlite_job_test.go:48,50
  internal/repository/sqlite/sqlite_job_test.go:51,53
  internal/repository/sqlite/sqlite_job_test.go:57,59
  internal/repository/sqlite/sqlite_listing_test.go:226,228
  internal/repository/sqlite/sqlite_listing_test.go:449,451
  internal/repository/sqlite/sqlite_user_test.go:91,93
  internal/repository/sqlite/sqlite_user_test.go:100,102
found 6 clones:
  internal/agent/verify.go:525,525
  internal/agent/verify_test.go:244,244
  internal/agent/verify_test.go:267,267
  internal/agent/verify_test.go:289,289
  internal/repository/cached/cached_test.go:138,138
  internal/repository/cached/cached_test.go:158,158
found 4 clones:
  cmd/security-audit/audit_test.go:75,80
  internal/domain/listing_validation_test.go:213,218
  internal/domain/listing_validation_test.go:275,280
  internal/ui/renderer_test.go:188,193
found 3 clones:
  internal/repository/cached/cached_test.go:207,209
  internal/repository/sqlite/sqlite_listing_test.go:155,157
  internal/repository/sqlite/sqlite_listing_test.go:158,160
found 3 clones:
  internal/ui/renderer.go:136,138
  internal/ui/renderer.go:142,144
  internal/ui/renderer.go:147,149
found 2 clones:
  internal/module/admin/admin_bulk_integration_test.go:125,125
  internal/service/csv.go:142,142
found 3 clones:
  cmd/harness/commands/root.go:53,55
  cmd/harness/commands/root.go:148,150
  cmd/harness/commands/update_coverage.go:45,47
found 2 clones:
  cmd/harness/commands/chaos.go:82,85
  internal/agent/security.go:116,116
found 12 clones:
  internal/handler/mock_repository_test.go:28,28
  internal/handler/mock_repository_test.go:68,68
  internal/handler/mock_repository_test.go:72,72
  internal/handler/mock_repository_test.go:108,108
  internal/module/admin/mock_repository_test.go:28,28
  internal/module/admin/mock_repository_test.go:68,68
  internal/module/admin/mock_repository_test.go:72,72
  internal/module/admin/mock_repository_test.go:108,108
  internal/repository/sqlite/sqlite_category.go:89,89
  internal/repository/sqlite/sqlite_listing_read.go:131,131
  internal/repository/sqlite/sqlite_user.go:57,57
  internal/repository/sqlite/sqlite_user.go:69,69
found 2 clones:
  internal/handler/ui_regression_admin_test.go:155,157
  internal/handler/ui_regression_admin_test.go:159,161
found 7 clones:
  internal/repository/sqlite/feedback_test.go:156,158
  internal/repository/sqlite/feedback_test.go:159,161
  internal/repository/sqlite/sqlite_listing_test.go:122,124
  internal/repository/sqlite/sqlite_listing_test.go:136,138
  internal/repository/sqlite/sqlite_listing_test.go:276,278
  internal/repository/sqlite/sqlite_user_test.go:44,46
  internal/service/csv_test.go:167,169
found 6 clones:
  internal/agent/security.go:237,237
  internal/agent/security.go:302,302
  internal/agent/security.go:341,341
  internal/agent/security.go:396,396
  internal/agent/security.go:448,448
  internal/agent/security.go:527,527
found 3 clones:
  internal/middleware/security_test.go:19,21
  internal/module/auth/middleware_test.go:31,33
  internal/module/auth/middleware_test.go:56,58
found 4 clones:
  internal/module/listing/listing_edge_cases_test.go:43,43
  internal/module/listing/listing_edge_cases_test.go:113,113
  internal/module/listing/listing_upload_test.go:59,59
  internal/module/listing/listing_upload_test.go:106,106
found 7 clones:
  internal/agent/drift_test.go:35,37
  internal/agent/drift_test.go:59,61
  internal/agent/drift_test.go:137,139
  internal/agent/drift_test.go:172,174
  internal/agent/progress_test.go:43,45
  internal/repository/cached/cached_test.go:106,108
  internal/repository/cached/cached_test.go:192,194
found 3 clones:
  internal/agent/security.go:221,227
  internal/agent/security.go:592,598
  internal/agent/security.go:612,618
found 2 clones:
  internal/module/listing/pagination_test.go:82,85
  internal/repository/cached/cached_test.go:26,26
found 2 clones:
  internal/repository/sqlite/sqlite_listing_read.go:230,237
  internal/repository/sqlite/sqlite_stats.go:19,26
found 2 clones:
  internal/module/listing/listing_mutations.go:90,92
  internal/module/listing/listing_mutations.go:143,145
found 2 clones:
  internal/repository/sqlite/sqlite_listing_write.go:58,58
  internal/repository/sqlite/sqlite_listing_write.go:121,121
found 2 clones:
  internal/module/admin/admin_export_test.go:54,54
  internal/module/admin/admin_export_test.go:55,55
found 9 clones:
  internal/agent/verify_apispec_test.go:164,187
  internal/agent/verify_apispec_test.go:214,233
  internal/agent/verify_apispec_test.go:262,270
  internal/agent/verify_apispec_test.go:290,298
  internal/agent/verify_apispec_test.go:318,326
  internal/agent/verify_apispec_test.go:372,386
  internal/agent/verify_apispec_test.go:488,508
  internal/agent/verify_apispec_test.go:538,546
  internal/agent/verify_apispec_test.go:571,588
found 5 clones:
  internal/middleware/ratelimit_test.go:22,24
  internal/middleware/ratelimit_test.go:42,44
  internal/middleware/ratelimit_test.go:73,75
  internal/middleware/ratelimit_test.go:92,94
  internal/middleware/ratelimit_test.go:126,128
found 2 clones:
  internal/module/admin/admin_actions_test.go:44,44
  internal/module/admin/admin_bulk_integration_test.go:27,27
found 71 clones:
  internal/seeder/seeder.go:59,59
  internal/seeder/seeder.go:60,60
  internal/seeder/seeder.go:61,61
  internal/seeder/seeder.go:62,62
  internal/seeder/seeder.go:63,63
  internal/seeder/seeder.go:64,64
  internal/seeder/seeder.go:65,65
  internal/seeder/seeder.go:66,66
  internal/seeder/seeder.go:67,67
  internal/seeder/seeder.go:68,68
  internal/seeder/seeder.go:72,72
  internal/seeder/seeder.go:73,73
  internal/seeder/seeder.go:74,74
  internal/seeder/seeder.go:75,75
  internal/seeder/seeder.go:76,76
  internal/seeder/seeder.go:77,77
  internal/seeder/seeder.go:78,78
  internal/seeder/seeder.go:79,79
  internal/seeder/seeder.go:80,80
  internal/seeder/seeder.go:81,81
  internal/seeder/seeder.go:85,85
  internal/seeder/seeder.go:86,86
  internal/seeder/seeder.go:87,87
  internal/seeder/seeder.go:88,88
  internal/seeder/seeder.go:89,89
  internal/seeder/seeder.go:90,90
  internal/seeder/seeder.go:91,91
  internal/seeder/seeder.go:92,92
  internal/seeder/seeder.go:93,93
  internal/seeder/seeder.go:94,94
  internal/seeder/seeder.go:95,95
  internal/seeder/seeder.go:99,99
  internal/seeder/seeder.go:100,100
  internal/seeder/seeder.go:101,101
  internal/seeder/seeder.go:102,102
  internal/seeder/seeder.go:103,103
  internal/seeder/seeder.go:104,104
  internal/seeder/seeder.go:105,105
  internal/seeder/seeder.go:106,106
  internal/seeder/seeder.go:107,107
  internal/seeder/seeder.go:108,108
  internal/seeder/seeder.go:112,112
  internal/seeder/seeder.go:113,113
  internal/seeder/seeder.go:114,114
  internal/seeder/seeder.go:115,115
  internal/seeder/seeder.go:116,116
  internal/seeder/seeder.go:117,117
  internal/seeder/seeder.go:118,118
  internal/seeder/seeder.go:119,119
  internal/seeder/seeder.go:120,120
  internal/seeder/seeder.go:121,121
  internal/seeder/seeder.go:125,125
  internal/seeder/seeder.go:126,126
  internal/seeder/seeder.go:127,127
  internal/seeder/seeder.go:128,128
  internal/seeder/seeder.go:129,129
  internal/seeder/seeder.go:130,130
  internal/seeder/seeder.go:131,131
  internal/seeder/seeder.go:132,132
  internal/seeder/seeder.go:133,133
  internal/seeder/seeder.go:134,134
  internal/seeder/seeder.go:138,138
  internal/seeder/seeder.go:139,139
  internal/seeder/seeder.go:140,140
  internal/seeder/seeder.go:141,141
  internal/seeder/seeder.go:142,142
  internal/seeder/seeder.go:143,143
  internal/seeder/seeder.go:144,144
  internal/seeder/seeder.go:145,145
  internal/seeder/seeder.go:146,146
  internal/seeder/seeder.go:147,147
found 2 clones:
  cmd/listing.go:142,144
  cmd/listing.go:145,147
found 5 clones:
  internal/repository/sqlite/sqlite_category_test.go:222,222
  internal/repository/sqlite/sqlite_category_test.go:223,223
  internal/repository/sqlite/sqlite_category_test.go:224,224
  internal/repository/sqlite/sqlite_category_test.go:225,225
  internal/repository/sqlite/sqlite_category_test.go:226,226
found 15 clones:
  internal/module/admin/admin_delete_test.go:27,27
  internal/module/admin/admin_delete_test.go:49,49
  internal/module/admin/admin_delete_test.go:115,115
  internal/module/admin/admin_delete_test.go:129,129
  internal/module/auth/handler_login_test.go:97,102
  internal/module/auth/handler_login_test.go:296,296
  internal/module/auth/handler_register_test.go:65,70
  internal/module/auth/handler_register_test.go:111,116
  internal/module/auth/handler_register_test.go:154,159
  internal/module/auth/handler_register_test.go:204,209
  internal/module/auth/middleware_test.go:79,79
  internal/repository/sqlite/sqlite_listing_test.go:216,221
  internal/repository/sqlite/sqlite_listing_test.go:319,319
  internal/repository/sqlite/sqlite_listing_test.go:320,320
  internal/repository/sqlite/sqlite_user_test.go:56,56
found 2 clones:
  internal/module/admin/admin_bulk.go:71,74
  internal/module/admin/admin_delete.go:69,72
found 17 clones:
  cmd/server.go:49,49
  internal/handler/mock_repository_test.go:20,20
  internal/handler/mock_repository_test.go:64,64
  internal/handler/mock_repository_test.go:76,76
  internal/handler/mock_repository_test.go:112,112
  internal/handler/mock_repository_test.go:116,116
  internal/module/admin/mock_repository_test.go:20,20
  internal/module/admin/mock_repository_test.go:64,64
  internal/module/admin/mock_repository_test.go:76,76
  internal/module/admin/mock_repository_test.go:112,112
  internal/module/admin/mock_repository_test.go:116,116
  internal/repository/sqlite/feedback.go:10,10
  internal/repository/sqlite/sqlite_category.go:11,11
  internal/repository/sqlite/sqlite_category.go:70,70
  internal/repository/sqlite/sqlite_claim.go:12,12
  internal/repository/sqlite/sqlite_listing_write.go:13,13
  internal/repository/sqlite/sqlite_user.go:23,23
found 4 clones:
  internal/agent/cost_test.go:16,16
  internal/agent/cost_test.go:27,27
  internal/ui/renderer_test.go:319,319
  internal/ui/renderer_test.go:323,323
found 14 clones:
  cmd/admin.go:30,33
  cmd/admin.go:53,56
  cmd/admin.go:76,79
  cmd/admin.go:102,105
  cmd/admin.go:137,140
  cmd/admin.go:176,179
  cmd/category.go:43,46
  cmd/category.go:59,62
  cmd/listing.go:103,106
  cmd/listing_backfill.go:32,35
  cmd/listing_read.go:26,29
  cmd/listing_read.go:61,64
  cmd/listing_update.go:21,24
  cmd/serve.go:77,80
found 6 clones:
  internal/module/auth/handler_login_test.go:97,102
  internal/module/auth/handler_login_test.go:296,296
  internal/module/auth/handler_register_test.go:65,70
  internal/module/auth/handler_register_test.go:111,116
  internal/module/auth/handler_register_test.go:154,159
  internal/module/auth/handler_register_test.go:204,209
found 4 clones:
  internal/agent/security.go:264,273
  internal/agent/security.go:355,364
  internal/agent/security.go:374,383
  internal/agent/security.go:425,436
found 2 clones:
  internal/ui/renderer_test.go:346,347
  internal/ui/renderer_test.go:358,359
found 5 clones:
  internal/domain/repository.go:15,15
  internal/domain/repository.go:53,53
  internal/domain/repository.go:66,66
  internal/domain/repository.go:73,73
  internal/domain/repository.go:74,74
found 5 clones:
  internal/repository/sqlite/repro_test.go:19,19
  internal/repository/sqlite/repro_test.go:20,20
  internal/repository/sqlite/sqlite_category_test.go:266,266
  internal/repository/sqlite/sqlite_category_test.go:267,267
  internal/repository/sqlite/sqlite_category_test.go:268,268
found 2 clones:
  internal/module/listing/listing_form_test.go:10,21
  internal/module/listing/listing_form_test.go:44,55
found 3 clones:
  internal/handler/ui_regression_admin_test.go:32,42
  internal/module/admin/admin_ui_integration_test.go:94,106
  internal/module/admin/admin_ui_integration_test.go:132,144
found 2 clones:
  cmd/security-audit/audit_test.go:20,24
  cmd/security-audit/audit_test.go:301,305
found 3 clones:
  internal/handler/ui_regression_admin_test.go:36,67
  internal/handler/ui_regression_admin_test.go:60,91
  internal/handler/ui_regression_admin_test.go:84,114
found 2 clones:
  internal/module/listing/ui_regression_home_test.go:87,89
  internal/module/listing/ui_regression_home_test.go:90,92
found 3 clones:
  internal/domain/listing_validation_test.go:416,429
  internal/domain/listing_validation_test.go:423,436
  internal/domain/listing_validation_test.go:430,443
found 4 clones:
  cmd/listing.go:138,138
  cmd/listing.go:143,143
  cmd/listing.go:146,146
  cmd/listing.go:151,151
found 3 clones:
  internal/repository/sqlite/sqlite_listing_test.go:85,88
  internal/repository/sqlite/sqlite_listing_test.go:91,94
  internal/repository/sqlite/sqlite_listing_test.go:97,100
found 5 clones:
  internal/module/auth/handler_login_test.go:162,162
  internal/module/auth/handler_login_test.go:187,187
  internal/module/auth/handler_login_test.go:212,212
  internal/module/auth/handler_logout_test.go:29,29
  internal/module/auth/handler_logout_test.go:48,48
found 4 clones:
  internal/module/auth/handler.go:260,260
  internal/module/auth/middleware.go:27,27
  internal/module/listing/listing.go:179,179
  internal/module/listing/listing.go:231,231
found 2 clones:
  cmd/server_public_test.go:20,26
  cmd/server_public_test.go:27,33
found 2 clones:
  cmd/security-audit/audit_test.go:149,156
  cmd/security-audit/audit_test.go:203,210
found 3 clones:
  internal/handler/mock_repository_test.go:128,128
  internal/module/admin/mock_repository_test.go:128,128
  internal/repository/sqlite/sqlite_claim.go:88,88
found 2 clones:
  internal/ui/renderer_test.go:283,285
  internal/ui/renderer_test.go:289,291
found 3 clones:
  internal/ui/renderer.go:65,65
  internal/ui/renderer.go:66,66
  internal/ui/renderer.go:67,67
found 5 clones:
  internal/module/admin/admin_category_test.go:49,49
  internal/module/admin/admin_category_test.go:67,67
  internal/module/admin/admin_category_test.go:92,92
  internal/module/admin/admin_middleware_test.go:31,31
  internal/module/listing/listing_edit_reproduction_test.go:142,142
found 2 clones:
  internal/agent/verify_apispec_test.go:119,125
  internal/agent/verify_test.go:224,230
found 2 clones:
  internal/seeder/seeder.go:86,135
  internal/seeder/seeder.go:99,148
found 2 clones:
  internal/agent/ast.go:32,37
  internal/module/listing/ui_regression_home_test.go:70,75
found 3 clones:
  cmd/benchmark.go:24,27
  cmd/seed.go:20,23
  cmd/stress.go:24,27
found 2 clones:
  internal/service/image_test.go:296,301
  internal/service/image_test.go:303,308
found 4 clones:
  internal/ui/renderer_test.go:113,118
  internal/ui/renderer_test.go:124,129
  internal/ui/renderer_test.go:135,140
  internal/ui/renderer_test.go:334,338
found 3 clones:
  cmd/harness/commands/chaos_test.go:30,32
  cmd/security-audit/audit_test.go:261,263
  internal/middleware/security_test.go:23,25
found 13 clones:
  internal/module/admin/admin_actions_test.go:68,68
  internal/module/admin/admin_actions_test.go:69,69
  internal/module/admin/admin_actions_test.go:70,70
  internal/module/admin/admin_actions_test.go:71,71
  internal/module/admin/admin_actions_test.go:80,80
  internal/module/admin/admin_actions_test.go:81,81
  internal/module/admin/admin_actions_test.go:82,82
  internal/module/admin/admin_actions_test.go:91,91
  internal/module/admin/admin_actions_test.go:92,92
  internal/module/admin/admin_actions_test.go:93,93
  internal/module/admin/admin_actions_test.go:94,94
  internal/repository/sqlite/sqlite_listing_test.go:349,349
  internal/repository/sqlite/sqlite_listing_test.go:350,350
found 2 clones:
  internal/module/auth/handler.go:256,263
  internal/module/listing/listing.go:227,234
found 4 clones:
  internal/handler/mock_repository_test.go:32,34
  internal/handler/mock_repository_test.go:56,58
  internal/module/admin/mock_repository_test.go:32,34
  internal/module/admin/mock_repository_test.go:56,58
found 5 clones:
  internal/module/listing/listing_featured_test.go:43,54
  internal/module/listing/listing_featured_test.go:98,109
  internal/module/listing/listing_featured_test.go:143,154
  internal/module/listing/listing_featured_test.go:182,193
  internal/module/listing/listing_featured_test.go:222,233
found 3 clones:
  internal/repository/cached/cached_test.go:49,54
  internal/repository/sqlite/feedback_test.go:189,194
  internal/repository/sqlite/feedback_test.go:192,197
found 2 clones:
  internal/module/auth/test_helpers_test.go:41,44
  internal/module/auth/test_helpers_test.go:49,52
found 2 clones:
  cmd/harness/commands/init.go:29,38
  internal/agent/state_test.go:68,77
found 3 clones:
  internal/module/listing/listing_delete_test.go:23,30
  internal/module/listing/listing_delete_test.go:44,51
  internal/module/listing/listing_update_test.go:24,31
found 2 clones:
  internal/agent/verify_test.go:303,311
  internal/agent/verify_test.go:333,341
found 3 clones:
  internal/repository/sqlite/repro_test.go:22,26
  internal/repository/sqlite/sqlite_category_test.go:96,100
  internal/repository/sqlite/sqlite_category_test.go:270,274
found 5 clones:
  internal/repository/sqlite/sqlite.go:100,102
  internal/repository/sqlite/sqlite.go:114,116
  internal/repository/sqlite/sqlite.go:127,129
  internal/repository/sqlite/sqlite.go:201,203
  internal/repository/sqlite/sqlite.go:217,219
found 2 clones:
  internal/agent/security_ast_test.go:17,24
  internal/agent/security_ast_test.go:41,49
found 2 clones:
  internal/agent/state_test.go:52,54
  internal/agent/state_test.go:55,57
found 4 clones:
  cmd/serve_test.go:19,42
  cmd/serve_test.go:27,50
  cmd/serve_test.go:35,58
  cmd/serve_test.go:43,66
found 4 clones:
  internal/module/listing/listing_read_test.go:37,37
  internal/module/listing/listing_service_test.go:31,31
  internal/module/listing/listing_service_test.go:79,79
  internal/module/listing/listing_service_test.go:91,91
found 2 clones:
  internal/ui/renderer_test.go:326,329
  internal/ui/renderer_test.go:359,361
found 6 clones:
  cmd/cli_json_test.go:27,29
  cmd/cli_json_test.go:44,46
  cmd/cli_json_test.go:68,70
  cmd/harness/commands/chaos_test.go:37,39
  cmd/harness/commands/chaos_test.go:72,74
  cmd/harness/commands/chaos_test.go:88,90
found 2 clones:
  internal/agent/cost.go:27,36
  internal/agent/drift.go:174,183
found 2 clones:
  cmd/harness/commands/commands_test.go:26,32
  cmd/harness/commands/commands_test.go:34,40
found 2 clones:
  internal/agent/coverage_test.go:39,43
  internal/agent/coverage_test.go:45,49
found 2 clones:
  internal/service/csv_test.go:124,129
  internal/service/csv_test.go:132,137
found 2 clones:
  internal/service/image_test.go:229,229
  internal/service/image_test.go:243,243
found 4 clones:
  internal/repository/sqlite/sqlite_claim.go:99,101
  internal/repository/sqlite/sqlite_listing_read.go:147,149
  internal/repository/sqlite/sqlite_user.go:62,64
  internal/repository/sqlite/sqlite_user.go:74,76
found 2 clones:
  internal/repository/cached/cached_test.go:46,51
  internal/repository/cached/cached_test.go:69,74
found 3 clones:
  internal/module/admin/admin_category_test.go:48,49
  internal/module/admin/admin_category_test.go:66,67
  internal/module/admin/admin_category_test.go:91,92
found 3 clones:
  internal/repository/sqlite/sqlite_listing_test.go:257,260
  internal/repository/sqlite/sqlite_listing_test.go:263,266
  internal/repository/sqlite/sqlite_listing_test.go:290,293
found 3 clones:
  internal/domain/listing_validation_test.go:20,25
  internal/domain/listing_validation_test.go:26,31
  internal/domain/listing_validation_test.go:38,43
found 4 clones:
  internal/domain/repository.go:9,9
  internal/handler/mock_repository_test.go:24,24
  internal/module/admin/mock_repository_test.go:24,24
  internal/repository/sqlite/sqlite_listing_read.go:54,54
found 8 clones:
  internal/handler/mock_repository_test.go:48,48
  internal/handler/mock_repository_test.go:84,84
  internal/module/admin/mock_repository_test.go:48,48
  internal/module/admin/mock_repository_test.go:84,84
  internal/repository/cached/cached.go:32,32
  internal/repository/cached/cached_test.go:21,21
  internal/repository/sqlite/sqlite_listing_read.go:214,214
  internal/repository/sqlite/sqlite_stats.go:10,10
found 2 clones:
  internal/module/admin/admin_claims.go:11,20
  internal/module/admin/admin_claims.go:23,32
found 2 clones:
  internal/handler/feedback_integration_test.go:120,143
  internal/handler/feedback_integration_test.go:145,168
found 2 clones:
  internal/handler/ui_helpers_test.go:47,55
  internal/module/listing/listing_helpers_test.go:69,77
found 2 clones:
  internal/agent/ast.go:180,186
  internal/agent/drift.go:107,113
found 2 clones:
  internal/module/admin/admin_routes_test.go:14,18
  internal/module/admin/admin_routes_test.go:20,24
found 4 clones:
  internal/handler/mock_repository_test.go:76,82
  internal/handler/mock_repository_test.go:116,122
  internal/module/admin/mock_repository_test.go:76,82
  internal/module/admin/mock_repository_test.go:116,122
found 2 clones:
  cmd/cli_json_test.go:57,59
  cmd/cli_json_test.go:80,82
found 2 clones:
  internal/module/listing/listing_upload_test.go:59,60
  internal/module/listing/listing_upload_test.go:106,107
found 9 clones:
  internal/agent/cost_test.go:38,38
  internal/agent/cost_test.go:44,44
  internal/agent/cost_test.go:50,50
  internal/agent/cost_test.go:54,54
  internal/ui/renderer_test.go:326,326
  internal/ui/renderer_test.go:346,346
  internal/ui/renderer_test.go:347,347
  internal/ui/renderer_test.go:358,358
  internal/ui/renderer_test.go:359,359
found 3 clones:
  internal/repository/sqlite/sqlite_featured_test.go:17,25
  internal/repository/sqlite/sqlite_featured_test.go:26,34
  internal/repository/sqlite/sqlite_featured_test.go:35,43
found 7 clones:
  cmd/harness/commands/verify_test.go:51,51
  internal/agent/verify_apispec_test.go:81,81
  internal/agent/verify_apispec_test.go:110,110
  internal/agent/verify_apispec_test.go:130,130
  internal/agent/verify_test.go:17,17
  internal/agent/verify_test.go:157,157
  internal/agent/verify_test.go:529,529
found 10 clones:
  internal/handler/mock_repository_test.go:20,22
  internal/handler/mock_repository_test.go:64,66
  internal/handler/mock_repository_test.go:76,78
  internal/handler/mock_repository_test.go:112,114
  internal/handler/mock_repository_test.go:116,118
  internal/module/admin/mock_repository_test.go:20,22
  internal/module/admin/mock_repository_test.go:64,66
  internal/module/admin/mock_repository_test.go:76,78
  internal/module/admin/mock_repository_test.go:112,114
  internal/module/admin/mock_repository_test.go:116,118
found 2 clones:
  cmd/serve_test.go:19,58
  cmd/serve_test.go:27,66
found 2 clones:
  internal/service/image_test.go:127,138
  internal/service/image_test.go:150,161
found 3 clones:
  internal/agent/verify_apispec_test.go:29,36
  internal/agent/verify_apispec_test.go:34,41
  internal/agent/verify_apispec_test.go:39,46
found 2 clones:
  internal/module/auth/handler_register_test.go:93,134
  internal/module/auth/handler_register_test.go:136,176
found 2 clones:
  cmd/admin.go:14,20
  cmd/category.go:14,19
found 3 clones:
  internal/agent/verify_apispec_test.go:252,283
  internal/agent/verify_apispec_test.go:280,311
  internal/agent/verify_apispec_test.go:362,399
found 2 clones:
  cmd/listing_cmd_test.go:150,155
  internal/module/listing/listing_ui_image_test.go:17,21
found 5 clones:
  internal/handler/test_db_util_test.go:19,23
  internal/module/admin/admin_bulk_test.go:152,152
  internal/module/admin/admin_bulk_test.go:174,174
  internal/module/admin/admin_dashboard_test.go:23,23
  internal/module/admin/admin_dashboard_test.go:24,24
found 2 clones:
  internal/repository/sqlite/sqlite_category_test.go:92,93
  internal/repository/sqlite/sqlite_category_test.go:93,94
found 3 clones:
  cmd/server_admin_test.go:97,97
  cmd/server_user_test.go:81,81
  cmd/server_user_test.go:98,98
found 3 clones:
  internal/domain/listing.go:154,158
  internal/domain/listing.go:160,164
  internal/domain/listing.go:166,170
found 7 clones:
  internal/module/auth/handler_login_test.go:242,242
  internal/module/auth/handler_login_test.go:268,268
  internal/module/auth/handler_login_test.go:293,293
  internal/module/auth/handler_register_test.go:35,35
  internal/module/auth/handler_register_test.go:62,62
  internal/module/auth/handler_register_test.go:108,108
  internal/module/auth/handler_register_test.go:151,151
found 2 clones:
  cmd/harness/commands/set_phase_test.go:19,23
  internal/agent/state_test.go:241,245
found 2 clones:
  internal/agent/security.go:261,261
  internal/handler/ui_regression_admin_test.go:164,164
found 3 clones:
  internal/module/listing/cache_busting_integration_test.go:39,40
  internal/module/listing/listing_edge_cases_test.go:116,116
  internal/module/listing/listing_update_image_test.go:41,41
found 2 clones:
  internal/domain/listing_validation_test.go:416,436
  internal/domain/listing_validation_test.go:423,443
found 3 clones:
  cmd/serve_test.go:19,50
  cmd/serve_test.go:27,58
  cmd/serve_test.go:35,66
found 2 clones:
  internal/agent/security.go:539,539
  internal/agent/security.go:539,539
found 2 clones:
  internal/agent/verify_test.go:168,175
  internal/agent/verify_test.go:177,184
found 2 clones:
  internal/handler/feedback.go:62,62
  internal/module/admin/admin.go:147,147
found 2 clones:
  internal/module/auth/handler.go:223,223
  internal/module/auth/handler.go:311,311
found 3 clones:
  internal/domain/address.go:10,12
  internal/domain/address.go:14,16
  internal/ui/renderer.go:116,120
found 3 clones:
  internal/repository/sqlite/sqlite_listing_test.go:65,65
  internal/repository/sqlite/sqlite_listing_test.go:66,66
  internal/repository/sqlite/sqlite_listing_test.go:67,67
found 2 clones:
  internal/module/listing/listing_edge_cases_test.go:47,55
  internal/module/listing/listing_upload_test.go:74,82
found 2 clones:
  internal/module/listing/listing_mutations.go:25,30
  internal/module/listing/listing_mutations.go:78,83
found 3 clones:
  internal/agent/verify_apispec_test.go:242,254
  internal/agent/verify_apispec_test.go:329,341
  internal/agent/verify_apispec_test.go:516,530
found 5 clones:
  internal/repository/sqlite/feedback_test.go:39,45
  internal/repository/sqlite/feedback_test.go:79,85
  internal/repository/sqlite/feedback_test.go:170,170
  internal/repository/sqlite/feedback_test.go:171,171
  internal/repository/sqlite/feedback_test.go:172,172
found 4 clones:
  internal/repository/sqlite/sqlite_listing_test.go:118,124
  internal/repository/sqlite/sqlite_listing_test.go:132,138
  internal/repository/sqlite/sqlite_listing_test.go:273,278
  internal/repository/sqlite/sqlite_user_test.go:40,46
found 3 clones:
  internal/module/admin/admin_all_listings_test.go:88,88
  internal/module/admin/admin_all_listings_test.go:89,89
  internal/module/admin/admin_all_listings_test.go:90,90
found 7 clones:
  internal/handler/mock_repository_test.go:32,32
  internal/handler/mock_repository_test.go:56,56
  internal/module/admin/mock_repository_test.go:32,32
  internal/module/admin/mock_repository_test.go:56,56
  internal/repository/sqlite/sqlite_listing_read.go:153,153
  internal/repository/sqlite/sqlite_listing_read.go:261,261
  internal/repository/sqlite/sqlite_stats.go:30,30
found 2 clones:
  internal/repository/sqlite/feedback_test.go:79,89
  internal/repository/sqlite/feedback_test.go:172,176
found 11 clones:
  internal/module/admin/admin_category_test.go:22,22
  internal/module/admin/admin_category_test.go:48,48
  internal/module/admin/admin_category_test.go:66,66
  internal/module/admin/admin_category_test.go:91,91
  internal/module/admin/admin_delete_test.go:32,32
  internal/module/admin/admin_delete_test.go:92,92
  internal/module/admin/admin_delete_test.go:108,108
  internal/module/admin/admin_delete_test.go:135,135
  internal/module/admin/admin_featured_mock_test.go:24,24
  internal/module/admin/admin_login_test.go:105,105
  internal/module/listing/listing_event_test.go:41,41
found 3 clones:
  internal/repository/sqlite/sqlite_listing_read.go:56,60
  internal/repository/sqlite/sqlite_listing_read.go:133,137
  internal/repository/sqlite/sqlite_listing_read.go:216,220
found 2 clones:
  internal/domain/category_test.go:12,19
  internal/repository/sqlite/sqlite_category_test.go:277,284
found 2 clones:
  internal/module/admin/admin_listings.go:35,38
  internal/module/listing/listing.go:119,122
found 4 clones:
  internal/handler/mock_repository_test.go:60,60
  internal/handler/response.go:13,13
  internal/module/admin/mock_repository_test.go:60,60
  internal/repository/sqlite/sqlite_listing_write.go:267,267
found 6 clones:
  internal/service/csv.go:88,93
  internal/service/csv.go:91,96
  internal/service/csv.go:94,99
  internal/service/csv.go:97,102
  internal/service/csv.go:100,105
  internal/service/csv.go:103,108
found 11 clones:
  internal/module/auth/handler.go:209,211
  internal/module/auth/handler.go:256,258
  internal/module/auth/handler.go:261,263
  internal/module/auth/handler.go:266,268
  internal/module/auth/handler.go:305,308
  internal/module/listing/listing.go:208,210
  internal/module/listing/listing.go:232,234
  internal/module/listing/listing_mutations.go:33,35
  internal/module/listing/listing_mutations.go:58,60
  internal/module/listing/listing_mutations.go:65,67
  internal/module/listing/listing_mutations.go:121,123
found 2 clones:
  internal/agent/drift.go:240,244
  internal/agent/drift.go:245,249
found 6 clones:
  internal/module/listing/ui_regression_home_test.go:32,32
  internal/repository/sqlite/sqlite_listing_test.go:254,254
  internal/repository/sqlite/sqlite_listing_test.go:286,286
  internal/repository/sqlite/sqlite_listing_test.go:381,381
  internal/repository/sqlite/sqlite_listing_test.go:382,382
  internal/repository/sqlite/sqlite_listing_test.go:383,383
found 2 clones:
  cmd/security-audit/audit_test.go:234,239
  cmd/security-audit/audit_test.go:342,347
found 3 clones:
  internal/service/csv.go:88,102
  internal/service/csv.go:91,105
  internal/service/csv.go:94,108
found 3 clones:
  cmd/admin.go:22,43
  cmd/admin.go:45,66
  cmd/admin.go:168,189
found 5 clones:
  internal/seeder/seeder.go:58,69
  internal/seeder/seeder.go:71,82
  internal/seeder/seeder.go:98,109
  internal/seeder/seeder.go:111,122
  internal/seeder/seeder.go:124,135
found 8 clones:
  internal/handler/mock_repository_test.go:80,82
  internal/handler/mock_repository_test.go:96,98
  internal/handler/mock_repository_test.go:100,102
  internal/handler/mock_repository_test.go:120,122
  internal/module/admin/mock_repository_test.go:80,82
  internal/module/admin/mock_repository_test.go:96,98
  internal/module/admin/mock_repository_test.go:100,102
  internal/module/admin/mock_repository_test.go:120,122
found 7 clones:
  internal/domain/repository.go:21,21
  internal/domain/repository.go:22,22
  internal/domain/repository.go:35,35
  internal/domain/repository.go:45,45
  internal/domain/repository.go:52,52
  internal/domain/repository.go:65,65
  internal/domain/repository.go:81,81
found 2 clones:
  cmd/listing_test.go:124,124
  cmd/listing_test.go:128,128
found 5 clones:
  internal/handler/ui_regression_modals_test.go:26,28
  internal/handler/ui_regression_modals_test.go:30,32
  internal/handler/ui_regression_modals_test.go:34,36
  internal/handler/ui_regression_modals_test.go:38,40
  internal/handler/ui_regression_modals_test.go:42,44
found 2 clones:
  internal/handler/integration_ds_test.go:109,111
  internal/middleware/session_test.go:44,46
found 2 clones:
  internal/module/listing/listing_claim.go:30,30
  internal/module/listing/listing_mutations.go:155,155
found 3 clones:
  internal/domain/listing_job_test.go:66,66
  internal/repository/sqlite/feedback_test.go:126,126
  internal/service/background_test.go:19,19
found 7 clones:
  internal/module/admin/admin_actions_test.go:36,36
  internal/module/listing/listing_create_test.go:21,21
  internal/module/listing/listing_delete_test.go:20,20
  internal/module/listing/listing_form_integration_test.go:21,21
  internal/module/listing/listing_form_integration_test.go:23,23
  internal/module/listing/listing_update_test.go:21,21
  internal/module/listing/listing_update_test.go:73,73
found 9 clones:
  internal/common/page_handler_test.go:22,22
  internal/handler/real_template_helpers_test.go:19,19
  internal/handler/ui_helpers_test.go:19,19
  internal/module/admin/admin_helpers_test.go:15,15
  internal/module/admin/ui_helpers_test.go:19,19
  internal/module/auth/test_helpers_test.go:19,19
  internal/module/listing/listing_helpers_test.go:21,21
  internal/module/listing/ui_helpers_test.go:19,19
  internal/ui/renderer.go:157,157
found 2 clones:
  cmd/aglog/main_test.go:138,140
  cmd/aglog/main_test.go:144,146
found 2 clones:
  internal/agent/coverage.go:91,99
  internal/agent/state.go:63,71
found 2 clones:
  internal/repository/sqlite/sqlite_listing_test.go:349,349
  internal/repository/sqlite/sqlite_listing_test.go:350,350
found 10 clones:
  cmd/aglog/main.go:28,28
  cmd/harness/commands/chaos.go:25,25
  cmd/harness/commands/cost.go:15,15
  cmd/harness/commands/gate.go:16,16
  cmd/harness/commands/handoff.go:16,16
  cmd/harness/commands/init.go:15,15
  cmd/harness/commands/set_phase.go:14,14
  cmd/harness/commands/status.go:15,15
  cmd/harness/commands/update_coverage.go:17,17
  cmd/harness/commands/verify.go:18,18
found 3 clones:
  internal/repository/sqlite/sqlite.go:38,43
  internal/repository/sqlite/sqlite.go:41,46
  internal/repository/sqlite/sqlite.go:44,49
found 4 clones:
  internal/module/admin/admin_bulk_test.go:216,216
  internal/module/admin/admin_bulk_test.go:217,217
  internal/module/admin/admin_bulk_test.go:289,289
  internal/module/admin/admin_bulk_test.go:290,290
found 4 clones:
  cmd/admin.go:102,120
  cmd/admin.go:137,155
  cmd/category.go:59,77
  cmd/listing_read.go:26,44
found 2 clones:
  internal/module/listing/listing.go:232,239
  internal/module/listing/listing_mutations.go:65,72
found 3 clones:
  internal/repository/sqlite/repro_test.go:71,74
  internal/repository/sqlite/repro_test.go:77,80
  internal/repository/sqlite/repro_test.go:112,115
found 2 clones:
  internal/service/csv.go:88,105
  internal/service/csv.go:91,108
found 2 clones:
  internal/module/listing/listing_service_test.go:30,31
  internal/module/listing/listing_service_test.go:90,91
found 2 clones:
  internal/handler/ui_regression_admin_test.go:28,28
  internal/handler/ui_regression_admin_test.go:131,131
found 2 clones:
  internal/domain/repository.go:52,53
  internal/domain/repository.go:65,66
found 2 clones:
  internal/domain/listing_event_test.go:20,26
  internal/domain/listing_event_test.go:41,47
found 2 clones:
  internal/module/admin/ui_helpers_test.go:1,63
  internal/module/listing/ui_helpers_test.go:1,63
found 2 clones:
  internal/agent/security.go:24,30
  internal/service/image.go:25,31
found 8 clones:
  internal/agent/ast_test.go:38,40
  internal/agent/ast_test.go:41,43
  internal/seeder/category_seeder_test.go:37,39
  internal/seeder/category_seeder_test.go:90,92
  internal/ui/renderer_test.go:59,61
  internal/ui/renderer_test.go:98,100
  internal/ui/renderer_test.go:148,150
  internal/ui/renderer_test.go:176,178
found 3 clones:
  internal/ui/renderer_test.go:124,126
  internal/ui/renderer_test.go:163,165
  internal/ui/renderer_test.go:206,208
found 4 clones:
  internal/agent/ast.go:182,185
  internal/agent/drift.go:109,112
  internal/agent/drift.go:128,131
  internal/util/slices.go:13,16
found 4 clones:
  internal/module/listing/listing_form_integration_test.go:34,34
  internal/module/listing/listing_form_integration_test.go:36,36
  internal/module/listing/listing_form_integration_test.go:49,49
  internal/module/listing/listing_form_integration_test.go:63,63
found 3 clones:
  internal/repository/sqlite/sqlite_category.go:11,27
  internal/repository/sqlite/sqlite_category.go:70,86
  internal/repository/sqlite/sqlite_claim.go:12,23
found 22 clones:
  cmd/admin.go:30,33
  cmd/admin.go:36,39
  cmd/admin.go:53,56
  cmd/admin.go:59,62
  cmd/admin.go:76,79
  cmd/admin.go:82,85
  cmd/admin.go:102,105
  cmd/admin.go:137,140
  cmd/admin.go:176,179
  cmd/admin.go:182,185
  cmd/category.go:43,46
  cmd/category.go:59,62
  cmd/listing.go:103,106
  cmd/listing_backfill.go:32,35
  cmd/listing_create.go:74,77
  cmd/listing_delete.go:19,22
  cmd/listing_read.go:26,29
  cmd/listing_read.go:61,64
  cmd/listing_update.go:21,24
  cmd/listing_update.go:89,92
  cmd/serve.go:77,80
  cmd/stress.go:37,40
found 4 clones:
  internal/module/auth/handler.go:37,37
  internal/module/auth/handler.go:100,100
  internal/module/auth/handler.go:152,152
  internal/module/auth/test_helpers_test.go:39,39
found 2 clones:
  internal/domain/listing_validation_test.go:311,315
  internal/domain/listing_validation_test.go:463,467
found 6 clones:
  cmd/admin.go:36,39
  cmd/admin.go:59,62
  cmd/admin.go:82,85
  cmd/admin.go:182,185
  cmd/listing_create.go:74,77
  cmd/listing_update.go:89,92
found 4 clones:
  internal/agent/verify_apispec_test.go:29,31
  internal/agent/verify_apispec_test.go:34,36
  internal/agent/verify_apispec_test.go:39,41
  internal/agent/verify_apispec_test.go:44,46
found 2 clones:
  internal/module/listing/listing_mutations.go:40,40
  internal/module/listing/listing_mutations.go:140,140
found 2 clones:
  internal/module/admin/admin_delete.go:20,22
  internal/module/admin/admin_delete.go:46,48
found 7 clones:
  cmd/admin.go:116,120
  cmd/admin.go:151,155
  cmd/category.go:73,77
  cmd/listing_create.go:79,83
  cmd/listing_read.go:40,44
  cmd/listing_read.go:66,70
  cmd/listing_update.go:94,98
found 3 clones:
  internal/agent/verify_apispec_test.go:628,630
  internal/agent/verify_apispec_test.go:631,633
  internal/agent/verify_apispec_test.go:634,636
found 2 clones:
  internal/agent/verify_test.go:232,276
  internal/agent/verify_test.go:255,298
found 3 clones:
  internal/handler/mock_repository_test.go:124,124
  internal/module/admin/mock_repository_test.go:124,124
  internal/repository/sqlite/sqlite_claim.go:51,51
found 2 clones:
  internal/agent/verify_test.go:65,82
  internal/agent/verify_test.go:117,134
found 3 clones:
  internal/service/image_test.go:118,118
  internal/service/image_test.go:185,185
  internal/service/image_test.go:204,204
found 7 clones:
  internal/module/auth/handler_login_test.go:105,106
  internal/module/auth/handler_login_test.go:298,299
  internal/module/auth/handler_register_test.go:40,41
  internal/module/auth/handler_register_test.go:81,82
  internal/module/auth/handler_register_test.go:127,128
  internal/module/auth/handler_register_test.go:170,171
  internal/module/auth/handler_register_test.go:211,212
found 3 clones:
  internal/module/admin/admin_listings.go:90,90
  internal/module/listing/listing.go:213,213
  internal/module/listing/listing_service.go:49,49
found 5 clones:
  internal/domain/repository.go:10,10
  internal/domain/repository.go:12,12
  internal/domain/repository.go:46,46
  internal/domain/repository.go:47,47
  internal/domain/repository.go:80,80
found 10 clones:
  internal/service/geocoding_test.go:22,37
  internal/service/geocoding_test.go:38,52
  internal/service/geocoding_test.go:59,66
  internal/service/geocoding_test.go:67,81
  internal/service/geocoding_test.go:82,96
  internal/service/geocoding_test.go:97,111
  internal/service/geocoding_test.go:112,119
  internal/service/geocoding_test.go:127,134
  internal/service/geocoding_test.go:135,142
  internal/service/geocoding_test.go:143,157
found 2 clones:
  cmd/listing_create.go:74,83
  cmd/listing_update.go:89,98
found 2 clones:
  internal/repository/sqlite/sqlite_listing_test.go:145,147
  internal/repository/sqlite/sqlite_listing_test.go:146,148
found 6 clones:
  cmd/serve_test.go:19,26
  cmd/serve_test.go:27,34
  cmd/serve_test.go:35,42
  cmd/serve_test.go:43,50
  cmd/serve_test.go:51,58
  cmd/serve_test.go:59,66
found 12 clones:
  internal/handler/mock_repository_test.go:80,80
  internal/handler/mock_repository_test.go:96,96
  internal/handler/mock_repository_test.go:100,100
  internal/handler/mock_repository_test.go:120,120
  internal/module/admin/mock_repository_test.go:80,80
  internal/module/admin/mock_repository_test.go:96,96
  internal/module/admin/mock_repository_test.go:100,100
  internal/module/admin/mock_repository_test.go:120,120
  internal/repository/sqlite/feedback.go:22,22
  internal/repository/sqlite/sqlite_claim.go:26,26
  internal/repository/sqlite/sqlite_stats.go:53,53
  internal/repository/sqlite/sqlite_stats.go:64,64
found 7 clones:
  cmd/listing_create.go:58,62
  cmd/listing_create.go:63,67
  cmd/listing_create.go:68,72
  cmd/listing_update.go:56,60
  cmd/listing_update.go:61,65
  cmd/listing_update.go:66,70
  cmd/listing_update.go:74,78
found 3 clones:
  internal/module/admin/admin_category_test.go:36,39
  internal/module/admin/admin_category_test.go:56,59
  internal/module/admin/admin_category_test.go:80,83
found 2 clones:
  internal/module/listing/listing_service.go:36,38
  internal/module/listing/listing_service.go:45,47
found 3 clones:
  internal/agent/verify_apispec_test.go:81,87
  internal/agent/verify_apispec_test.go:110,116
  internal/agent/verify_apispec_test.go:130,136
found 3 clones:
  internal/agent/security_test.go:16,18
  internal/util/fs_test.go:114,116
  internal/util/fs_test.go:196,198
found 6 clones:
  internal/module/admin/admin_ui_integration_test.go:28,33
  internal/module/admin/admin_ui_integration_test.go:34,39
  internal/repository/sqlite/sqlite_listing_test.go:287,287
  internal/repository/sqlite/sqlite_user_test.go:16,16
  internal/repository/sqlite/sqlite_user_test.go:17,17
  internal/repository/sqlite/sqlite_user_test.go:33,33
found 2 clones:
  internal/module/admin/admin.go:22,32
  internal/module/admin/admin.go:34,44
found 3 clones:
  internal/handler/feedback.go:16,18
  internal/module/auth/middleware.go:17,19
  internal/service/background.go:15,17
found 3 clones:
  cmd/listing_create.go:74,83
  cmd/listing_read.go:61,70
  cmd/listing_update.go:89,98
found 2 clones:
  internal/handler/ui_regression_admin_test.go:36,91
  internal/handler/ui_regression_admin_test.go:60,114
found 3 clones:
  internal/repository/sqlite/sqlite_listing_test.go:74,79
  internal/repository/sqlite/sqlite_listing_test.go:86,91
  internal/repository/sqlite/sqlite_listing_test.go:92,97
found 3 clones:
  internal/handler/mock_repository_test.go:24,24
  internal/module/admin/mock_repository_test.go:24,24
  internal/repository/sqlite/sqlite_listing_read.go:54,54
found 2 clones:
  internal/handler/integration_ds_test.go:50,52
  internal/handler/integration_ds_test.go:112,114
found 9 clones:
  internal/module/admin/admin_bulk_test.go:135,135
  internal/module/admin/admin_bulk_test.go:150,150
  internal/module/admin/admin_bulk_test.go:172,172
  internal/module/admin/admin_bulk_test.go:194,194
  internal/module/admin/admin_bulk_test.go:213,213
  internal/module/admin/admin_bulk_test.go:278,278
  internal/module/admin/admin_delete_test.go:50,50
  internal/module/admin/admin_delete_test.go:65,65
  internal/module/admin/admin_delete_test.go:77,77
found 18 clones:
  cmd/cli_json_test.go:33,35
  internal/agent/verify_apispec_test.go:203,205
  internal/agent/verify_apispec_test.go:249,251
  internal/agent/verify_apispec_test.go:252,254
  internal/agent/verify_apispec_test.go:280,282
  internal/agent/verify_apispec_test.go:308,310
  internal/agent/verify_apispec_test.go:336,338
  internal/agent/verify_apispec_test.go:339,341
  internal/agent/verify_apispec_test.go:362,364
  internal/agent/verify_apispec_test.go:396,398
  internal/agent/verify_apispec_test.go:434,436
  internal/agent/verify_apispec_test.go:477,479
  internal/agent/verify_apispec_test.go:524,526
  internal/agent/verify_apispec_test.go:528,530
  internal/agent/verify_apispec_test.go:559,561
  internal/agent/verify_apispec_test.go:603,605
  internal/agent/verify_truncation_test.go:26,28
  internal/agent/verify_truncation_test.go:32,34
found 2 clones:
  internal/domain/listing.go:178,178
  internal/domain/listing.go:226,226
found 3 clones:
  internal/module/auth/provider_test.go:19,24
  internal/module/auth/provider_test.go:25,30
  internal/module/auth/provider_test.go:43,48
found 2 clones:
  internal/module/listing/listing.go:227,229
  internal/module/listing/listing_mutations.go:115,117
found 2 clones:
  internal/seeder/stress_generator.go:97,97
  internal/seeder/stress_generator.go:102,102
found 7 clones:
  internal/common/page_handler_test.go:22,24
  internal/handler/real_template_helpers_test.go:19,21
  internal/handler/ui_helpers_test.go:19,21
  internal/module/admin/ui_helpers_test.go:19,21
  internal/module/auth/test_helpers_test.go:19,21
  internal/module/listing/listing_helpers_test.go:21,23
  internal/module/listing/ui_helpers_test.go:19,21
found 5 clones:
  internal/agent/ast_test.go:62,64
  internal/agent/verify_test.go:65,67
  internal/agent/verify_test.go:117,119
  internal/util/fs_test.go:57,59
  internal/util/fs_test.go:76,78
found 17 clones:
  internal/module/admin/admin_actions_test.go:20,20
  internal/module/admin/admin_actions_test.go:113,113
  internal/module/admin/admin_actions_test.go:142,142
  internal/module/admin/admin_actions_test.go:160,160
  internal/module/admin/admin_all_listings_test.go:73,73
  internal/module/admin/admin_all_listings_test.go:93,93
  internal/module/admin/admin_bulk_test.go:75,75
  internal/module/admin/admin_bulk_test.go:241,241
  internal/module/admin/admin_claims_test.go:26,26
  internal/module/admin/admin_claims_test.go:45,45
  internal/module/admin/admin_claims_test.go:61,61
  internal/module/admin/admin_claims_test.go:79,79
  internal/module/admin/admin_dashboard_test.go:41,41
  internal/module/admin/admin_dashboard_test.go:60,60
  internal/module/admin/admin_featured_mock_test.go:27,27
  internal/module/admin/admin_featured_mock_test.go:43,43
  internal/module/admin/admin_users_test.go:18,18
found 3 clones:
  internal/repository/sqlite/sqlite_category_test.go:92,92
  internal/repository/sqlite/sqlite_category_test.go:93,93
  internal/repository/sqlite/sqlite_category_test.go:94,94
found 2 clones:
  cmd/security-audit/main.go:98,100
  cmd/security-audit/main.go:131,133
found 4 clones:
  internal/module/listing/listing_create_test.go:56,56
  internal/module/listing/listing_edge_cases_test.go:86,86
  internal/module/listing/listing_form_integration_test.go:91,91
  internal/module/listing/listing_update_test.go:98,98
found 2 clones:
  internal/repository/sqlite/sqlite_category.go:54,54
  internal/repository/sqlite/sqlite_category.go:99,99
found 2 clones:
  internal/module/auth/provider_test.go:31,36
  internal/module/auth/provider_test.go:37,42

Found total 465 clone groups.
```
