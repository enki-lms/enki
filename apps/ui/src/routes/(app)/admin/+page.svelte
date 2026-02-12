<script lang="ts">
  import Tabs from "$lib/components/Tabs.svelte";
  import Card from "$lib/components/Card.svelte";
  import Button from "$lib/components/Button.svelte";
  import LogoPlaceHolder from "$lib/components/LogoPlaceHolder.svelte";
  import ProfileMenu from "$lib/components/ProfileMenu.svelte";
  import FileManager from "$lib/components/FileManager.svelte";
  import ProblemSets from "$lib/components/ProblemSets.svelte";
  import Lectures from "$lib/components/Lectures.svelte";
  import Courses from "$lib/components/Courses.svelte";
  import CourseMaterials from "$lib/components/CourseMaterials.svelte";
  import Students from "$lib/components/Students.svelte";

  export let data;

  type TabType =
    | "Problem Sets"
    | "Lectures"
    | "Exams"
    | "Courses"
    | "Students"
    | "Course Materials";

  const tabsList: TabType[] = [
    "Problem Sets",
    "Lectures",
    "Exams",
    "Courses",
    "Students",
    "Course Materials",
  ];
  let activeTab: TabType = "Problem Sets";
</script>

<div class="min-h-screen p-8 bg-gray-100">
  <div class="max-w-6xl mx-auto">
    <!-- Header with Avatar and Logo -->
    <div class="mb-6 flex items-center gap-6">
      <ProfileMenu
        width="80"
        height="80"
        name={data.user?.fullName ?? ""}
        email={data.user?.email ?? ""}
        role={data.user?.role ?? ""}
      />
      <LogoPlaceHolder />
    </div>

    <!-- Main Card -->
    <Card
      padding="p-0"
      shadow={false}
      class="bg-white rounded-xl border-gray-200"
    >
      <!-- Tabs -->
      <div class="px-6 pt-6">
        <Tabs tabs={tabsList} bind:active={activeTab} />
      </div>

      <!-- Tab Content Area -->
      <div class="px-6 pb-6 pt-8 min-h-[600px]">
        {#if activeTab === "Problem Sets"}
          <ProblemSets />
        {:else if activeTab === "Lectures"}
          <Lectures />
        {:else if activeTab === "Exams"}
          <div class="text-center py-12">
            <h3 class="text-lg font-medium text-gray-900">Exam Dashboard</h3>
            <p class="mt-1 text-sm text-gray-500 mb-6">
              Manage exam sessions, grading, and invigilation in the dedicated
              dashboard.
            </p>
            <Button href="/admin/exams" variant="primary"
              >Go to Exam Dashboard</Button
            >
          </div>
        {:else if activeTab === "Courses"}
          <Courses />
        {:else if activeTab === "Students"}
          <Students />
        {:else if activeTab === "Course Materials"}
          <CourseMaterials />
        {/if}
      </div>
    </Card>
  </div>
</div>
