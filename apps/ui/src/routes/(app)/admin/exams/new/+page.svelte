<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import {
    api,
    type Course,
    type ProblemGroup,
    type Student,
    type QuizGroup,
    type Enrollment,
  } from "$lib/api";
  import Button from "$lib/components/Button.svelte";
  import Card from "$lib/components/Card.svelte";

  let courses: Course[] = [];
  let selectedCourseId: number | null = null;

  let problemGroups: { id: number; name: string; type: string }[] = [];
  let selectedGroupId: number | null = null;
  let selectedGroupType: "comp_sci" | "quiz" = "comp_sci";

  let students: Enrollment[] = [];
  let selectedStudentIds: number[] = [];

  let durationMinutes = 60;

  let isLoading = true;
  let isSubmitting = false;
  let error: string | null = null;

  onMount(async () => {
    try {
      courses = await api.getCourses(); // Teacher's courses
    } catch (e) {
      error = (e as Error).message;
    } finally {
      isLoading = false;
    }
  });

  async function handleCourseChange() {
    if (!selectedCourseId) {
      problemGroups = [];
      students = [];
      return;
    }

    try {
      const [csGroups, quizGroups, courseStudents] = await Promise.all([
        api.getCourseProblemGroups(selectedCourseId),
        api.getCourseQuizGroups(selectedCourseId),
        api.getCourseEnrollments(selectedCourseId),
      ]);

      problemGroups = [
        ...csGroups.map((g) => ({ ...g, type: "comp_sci" })),
        ...quizGroups.map((g) => ({ ...g, type: "quiz" })),
      ];
      students = courseStudents;
      // Select all students by default
      selectedStudentIds = students.map((s) => s.user_id);
    } catch (e) {
      console.error(e);
      error = "Failed to load course data";
    }
  }

  function handleGroupChange() {
    const group = problemGroups.find((g) => g.id === selectedGroupId);
    if (group) {
      selectedGroupType = group.type as "comp_sci" | "quiz";
    }
  }

  async function handleSubmit() {
    if (!selectedGroupId || selectedStudentIds.length === 0) return;

    isSubmitting = true;
    error = null;

    try {
      await api.createExamSession({
        problem_group_type: selectedGroupType,
        problem_group_id: selectedGroupId,
        duration_minutes: durationMinutes,
        student_ids: selectedStudentIds,
      });
      goto("/admin/exams");
    } catch (e) {
      error = (e as Error).message;
    } finally {
      isSubmitting = false;
    }
  }

  function toggleStudent(id: number) {
    if (selectedStudentIds.includes(id)) {
      selectedStudentIds = selectedStudentIds.filter((sid) => sid !== id);
    } else {
      selectedStudentIds = [...selectedStudentIds, id];
    }
  }

  function toggleAllStudents() {
    if (selectedStudentIds.length === students.length) {
      selectedStudentIds = [];
    } else {
      selectedStudentIds = students.map((s) => s.user_id);
    }
  }
</script>

<div class="max-w-3xl mx-auto space-y-6">
  <div class="flex items-center gap-4">
    <Button href="/admin/exams" variant="secondary">Back</Button>
    <div>
      <h1 class="text-2xl font-bold text-gray-900">Create New Exam Session</h1>
      <p class="text-gray-600">Schedule an exam for your students.</p>
    </div>
  </div>

  {#if error}
    <div class="bg-red-50 text-red-700 p-4 rounded-lg">
      {error}
    </div>
  {/if}

  <Card>
    <div class="space-y-6">
      <!-- Course Selection -->
      <div>
        <label for="course" class="block text-sm font-medium text-gray-700"
          >Course</label
        >
        <select
          id="course"
          bind:value={selectedCourseId}
          on:change={handleCourseChange}
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
        >
          <option value={null}>Select a course</option>
          {#each courses as course}
            <option value={course.id}>{course.name}</option>
          {/each}
        </select>
      </div>

      {#if selectedCourseId}
        <!-- Problem Group Selection -->
        <div>
          <label for="group" class="block text-sm font-medium text-gray-700"
            >Problem Group / Quiz</label
          >
          <select
            id="group"
            bind:value={selectedGroupId}
            on:change={handleGroupChange}
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
          >
            <option value={null}>Select a problem group</option>
            {#each problemGroups as group}
              <option value={group.id}>{group.name} ({group.type})</option>
            {/each}
          </select>
        </div>

        <!-- Duration -->
        <div>
          <label for="duration" class="block text-sm font-medium text-gray-700"
            >Duration (Minutes)</label
          >
          <input
            type="number"
            id="duration"
            bind:value={durationMinutes}
            min="1"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
          />
        </div>

        <!-- Students Selection -->
        <div>
          <div class="flex justify-between items-center mb-2">
            <label class="block text-sm font-medium text-gray-700"
              >Students</label
            >
            <button
              type="button"
              class="text-sm text-indigo-600 hover:text-indigo-900"
              on:click={toggleAllStudents}
            >
              {selectedStudentIds.length === students.length
                ? "Deselect All"
                : "Select All"}
            </button>
          </div>
          <div class="border rounded-md max-h-60 overflow-y-auto p-2 space-y-2">
            {#each students as student}
              <div class="flex items-center">
                <input
                  type="checkbox"
                  id={`student-${student.user_id}`}
                  checked={selectedStudentIds.includes(student.user_id)}
                  on:change={() => toggleStudent(student.user_id)}
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label
                  for={`student-${student.user_id}`}
                  class="ml-2 block text-sm text-gray-900"
                >
                  {student.user?.full_name ||
                    student.user?.email ||
                    `User #${student.user_id}`}
                </label>
              </div>
            {/each}
            {#if students.length === 0}
              <p class="text-gray-500 text-sm italic">
                No students enrolled in this course.
              </p>
            {/if}
          </div>
          <p class="text-xs text-gray-500 mt-1">
            {selectedStudentIds.length} students selected
          </p>
        </div>
      {/if}

      <div class="flex justify-end pt-4">
        <Button
          on:click={handleSubmit}
          disabled={isSubmitting ||
            !selectedGroupId ||
            selectedStudentIds.length === 0}
          variant="primary"
        >
          {isSubmitting ? "Creating..." : "Create Exam Session"}
        </Button>
      </div>
    </div>
  </Card>
</div>
