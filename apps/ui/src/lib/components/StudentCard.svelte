<script lang="ts">
  interface StudentProp {
    id: string;
    name: string;
    email: string;
  }

  export let student: StudentProp;

  function getInitials(name: string): string {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2);
  }
</script>

<div
  class="student-card bg-white border border-gray-200 rounded-xl p-6 shadow-sm hover:shadow-md transition-shadow duration-200"
>
  <div class="flex items-center gap-4">
    <div
      class="w-14 h-14 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-semibold text-lg flex-shrink-0"
    >
      {getInitials(student.name)}
    </div>

    <div class="flex-1 min-w-0">
      <h3 class="font-semibold text-gray-900 text-lg truncate">
        {student.name}
      </h3>
      <p class="text-sm text-gray-500 truncate">{student.email}</p>
    </div>
  </div>
</div>

<style>
  .student-card {
    transition: transform 0.2s ease;
  }

  .student-card:hover {
    transform: translateY(-2px);
  }
</style>
