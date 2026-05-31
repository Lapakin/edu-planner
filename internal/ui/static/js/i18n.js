/* =========================================================
   EduPlanner — Internationalization (i18n)
   Supported languages: ua (default), en
   ========================================================= */

const TRANSLATIONS = {
  /* ── Ukrainian ──────────────────────────────────────── */
  ua: {
    nav: {
      home: "Головна",
      reference: "Довідники",
      processing: "Плани",
      schedule: "Розклад",
      academicSettings: "Академічний календар",
      users: "Користувачі",
      profile: "Особистий кабінет",
    },

    auth: {
      loginTitle: "Вхід до системи",
      loginSubtitle: "Введіть дані облікового запису для продовження.",
      emailLabel: "Електронна пошта",
      emailPlaceholder: "example@example.com",
      passwordLabel: "Пароль",
      loginBtn: "Увійти",
      logoutBtn: "Вийти",
      loginError: "Невірна електронна пошта або пароль.",
      accountDeactivated:
        "Ваш обліковий запис деактивовано. Зверніться до адміністратора.",
      setPasswordTitle: "Встановлення пароля",
      setPasswordBtn: "Встановити пароль",
      setPasswordOk: "Пароль успішно встановлено. Входимо…",
      confirmPassword: "Підтвердіть пароль",
      passwordMismatch: "Паролі не співпадають.",
      tokenInvalid: "Посилання недійсне або термін дії сплив.",
    },

    header: {
      academicYear: "Навчальний рік",
      semester: "Семестр",
      selectYear: "Оберіть рік",
      selectSemester: "Оберіть семестр",
      language: "Мова",
    },

    lang: {
      ua: "Українська",
      en: "English",
    },

    sidebar: {
      mainNav: "Навігація",
      academicSettings: "Академічний календар",
    },

    actions: {
      add: "Додати",
      save: "Зберегти",
      discard: "Скасувати",
      delete: "Видалити",
      deleteSelected: "Видалити вибрані",
      edit: "Редагувати",
      view: "Переглянути",
      close: "Закрити",
      cancel: "Відмінити",
      confirm: "Підтвердити",
      generate: "Згенерувати",
      activate: "Активувати",
      deactivate: "Деактивувати",
      optional: "необов'язково",
      copy: "Копіювати",
      back: "Назад",
      actions: "Дії",
    },

    messages: {
      saved: "Зміни успішно збережено.",
      saveError: "Помилка збереження даних.",
      deleted: "Записи успішно видалено.",
      deleteError: "Помилка видалення записів.",
      loadError: "Помилка завантаження даних.",
      confirmDelete: "Ви впевнені, що хочете видалити вибрані записи?",
      unsavedChanges:
        "Є незбережені зміни. Збережіть або скасуйте їх перед переходом.",
      noAcademicYear:
        "Навчальний рік не обрано. Будь ласка, оберіть або створіть навчальний рік.",
      noData: "Немає даних",
      fillRequired: "Заповніть усі обов'язкові поля.",
      forbidden: "Доступ заборонено.",
      copied: "Скопійовано!",
    },

    status: {
      active: "Активний",
      inactive: "Неактивний",
    },

    fields: {
      id: "ID",
      name: "Назва",
      shortName: "Скорочена назва",
      email: "Електронна пошта",
      role: "Роль",
      password: "Пароль",
      startDate: "Дата початку",
      endDate: "Дата завершення",
      periodStart: "Початок",
      periodEnd: "Кінець",
      isActive: "Активний",
      isDeleted: "Видалений",
      createdAt: "Створено",
      modifiedAt: "Змінено",
      academicYear: "Навчальний рік",
      semester: "Семестр",
      department: "Кафедра",
      specialty: "Спеціальність",
      cycleCommittee: "Циклова комісія",
      teacher: "Викладач",
      discipline: "Дисципліна",
      group: "Група",
      lessonNumber: "Номер пари",
      startTime: "Початок пари",
      endTime: "Кінець пари",
      roomType: "Тип аудиторії",
      specialtyCode: "Код спеціальності",
      isSplitting: "Поділ на підгрупи",
      isSubvention: "Субвенція",
      isContract: "Контракт",
      educationStart: "Початок навчання",
      educationEnd: "Кінець навчання",
      slotType: "Тип слоту",
      weekday: "День тижня",
      room: "Аудиторія",
      maxIdenticalLessonsPerDay: "Макс. однакових пар на день",
      maxStudyHoursPerDay: "Макс. навч. годин на день",
      maxTeacherHoursPerWeek: "Макс. годин викладача на тиждень",
      maxGroupLessonHoursPerWeek: "Макс. год. занять групи на тиждень",
      minGroupLessonsPerDay: "Мін. пар групи на день",
      maxGroupLessonsPerDay: "Макс. пар групи на день",
      maxTeacherLessonsPerDay: "Макс. пар викладача на день",
      noGapsInGroupSchedule: "Без вікон у розкладі групи",
      lessonsPerClass: "Тривалість заняття (акад. год.)",
      studyDaysMask: "Маска навчальних днів",
      maxConsecutiveTeacherLessons: "Макс. пар підряд для викладача",
      timePriority: "Пріоритет часу занять",
      allowFlowLessons: "Дозволити потокові заняття",
      lecturesHours: "Лекції (год.)",
      laboratoryHours: "Лабораторні (год.)",
      practicalHours: "Практичні (год.)",
      independentWork: "Самостійна робота (год.)",
      examHours: "Іспит (год.)",
      creditHours: "Залік (год.)",
      classroomWork: "Аудиторна робота (год.)",
      semesterNumber: "Номер семестру",
      total: "Усього",
      educationalComponent: "Навчальний компонент",
      studyPlan: "Навчальний план",
      workloadDistribution: "Розподіл навантаження",
      assignedHours: "Призначені години",
      roleType: "Тип ролі",
      userId: "ID користувача",
      fullName: "Повне ім'я",
      startDate: "Дата початку",
      endDate: "Дата завершення",
      group: "Група",
      semester: "Семестр",
      firstName: "Ім'я",
      lastName: "Прізвище",
      patronymic: "По батькові",
    },

    entities: {
      academicYears: "Навчальні роки",
      semesters: "Семестри",
      departments: "Кафедри",
      rooms: "Аудиторії",
      specialties: "Спеціальності",
      cycleCommittees: "Циклові комісії",
      teachers: "Викладачі",
      disciplines: "Дисципліни",
      groups: "Групи",
      groupSemesters: "Групи по семестрах",
      bellSchedules: "Розклад дзвінків",
      studyPlans: "Навчальні плани",
      workloadDistributions: "Розподіл навантаження",
      workloadAssignments: "Призначення навантаження",
      scheduleTemplateSettings: "Налаштування шаблону",
      scheduleRestrictions: "Обмеження розкладу",
      scheduleTemplates: "Шаблони розкладу",
      teacherSlotPreferences: "Переваги слотів викладача",
      cycleCommitteeLabRooms: "Лабораторні аудиторії комісій",
      users: "Користувачі",
    },

    pages: {
      homeTitle: "Головна",
      homeWelcome: "Ласкаво просимо до EduPlanner",
      referenceTitle: "Довідники",
      processingTitle: "Плани",
      scheduleTitle: "Розклад",
      academicSettingsTitle: "Академічний календар",
      profileTitle: "Особиста інформація",
      usersTitle: "Управління користувачами",
      setPasswordTitle: "Встановлення пароля",
    },

    roomTypes: {
      auditorium: "Аудиторія",
      laboratory: "Лабораторія",
    },

    lessonsPerClass: {
      one: "1 акад. година",
      two: "2 акад. години",
      three: "3 акад. години",
    },

    timePriority: {
      none: "Без пріоритету",
      morning: "Перша половина дня",
      afternoon: "Друга половина дня",
    },

    weekdays: {
      monday: "Понеділок",
      tuesday: "Вівторок",
      wednesday: "Середа",
      thursday: "Четвер",
      friday: "П'ятниця",
      saturday: "Субота",
    },

    slotTypes: {
      preferred: "Бажаний",
      forbidden: "Заборонений",
    },

    roles: {
      admin: "Адміністратор",
      dean: "Деканат",
      guest: "Гість",
      teacher: "Викладач",
    },

    units: {
      lessons: "пар",
      hours: "год",
      consecutive: "підряд",
    },

    weekdaysShort: {
      monday: "Пн",
      tuesday: "Вт",
      wednesday: "Ср",
      thursday: "Чт",
      friday: "Пт",
      saturday: "Сб",
      sunday: "Нд",
    },

    modal: {
      unsaved: "Є незбережені зміни",
      saved: "Усі зміни збережено",

      tabFormat: "Формат пари",
      tabCalendar: "Календар",
      tabDaily: "Денні ліміти",
      tabWeekly: "Тижневі ліміти",
      tabLoad: "Навантаження",
      tabTime: "Час і календар",
      tabSpace: "Аудиторії",
      tabOrg: "Організаційні",

      lessonDuration: "Тривалість заняття",
      lessonDurationDesc: "Кількість академічних годин в одній парі.",
      studyDays: "Навчальні дні",
      studyDaysDesc: "Дні тижня, у які формується розклад.",
      maxIdentical: "Макс. однакових пар на день",
      maxIdenticalDesc:
        "Скільки пар з однієї дисципліни може бути в один день.",
      maxHoursDay: "Макс. навчальних годин на день",
      maxHoursDayDesc:
        "Сумарна верхня межа аудиторних годин у розкладі групи за день.",
      maxTeacherWeek: "Макс. годин викладача на тиждень",
      maxTeacherWeekDesc:
        "Верхня межа аудиторного навантаження одного викладача за тиждень.",
      maxGroupWeek: "Макс. годин занять групи на тиждень",
      maxGroupWeekDesc:
        "Сумарна верхня межа годин у розкладі однієї групи за тиждень.",
      maxGroupDay: "Макс. пар на день для студентів",
      maxGroupDayDesc:
        "Верхня межа кількості пар у розкладі групи за один навчальний день.",
      maxTeacherDay: "Макс. пар на день для викладачів",
      maxTeacherDayDesc:
        "Скільки пар максимум може провести один викладач протягом дня.",
      consecutive: "Безперервна робота викладача",
      consecutiveDesc:
        "Максимум пар підряд без перерви для одного викладача (0 — без ліміту).",
      minGroupDay: "Мін. пар групи на день",
      minGroupDayDesc:
        "Якщо група має заняття у цей день — мінімальна кількість пар (0 — без ліміту).",
      gapPrevention: "Запобігання «вікнам»",
      gapPreventionDesc:
        "Мінімізувати порожні години між парами у студентів і викладачів.",
      gapOff: "Вимкнено",
      gapOn: "Виключити повністю",
      forbiddenHours: "Заборонені години викладачів",
      forbiddenHoursDesc: "Час, коли викладач не може проводити заняття.",
      teacherPrefs: "Побажання викладачів",
      teacherPrefsDesc:
        "Пріоритетні години, у які викладач хоче проводити заняття.",
      timePriority: "Пріоритет проведення занять",
      timePriorityDesc: "Бажана половина дня для розміщення пар у розкладі.",
      labRooms: "Лабораторні фахових дисциплін",
      labRoomsDesc:
        "Виберіть комп'ютерні класи, у яких можуть проводитися лабораторні роботи.",
      flowLessons: "Потокові заняття (спільні лекції)",
      flowLessonsDesc:
        "Дозволити одному викладачу вести пару одночасно для кількох груп.",
      flowOn: "Дозволено",
      flowOff: "Заборонено",
      configure: "Налаштувати",
      back: "Назад",
      noSettings: "Немає налаштувань",
      selectTeacher: "Оберіть викладача",
      lessonNum: "Пара №",
      chooseCycleCommittee: "Оберіть циклову комісію…",
      forbiddenTitle: "Заборонені години викладачів",
      preferredTitle: "Побажання викладачів",
    },

    setupModal: {
      title: "Налаштування навчального року",
      description:
        "Навчальний рік не виявлено. Будь ласка, заповніть дані для створення нового навчального року та першого семестру.",
      yearName: "Назва навчального року",
      yearStart: "Початок навчального року",
      yearEnd: "Кінець навчального року",
      semesterStart: "Початок семестру",
      semesterEnd: "Кінець семестру",
      createBtn: "Зберегти",
    },

    tabs: {
      /* Reference */
      academicYears: "Навчальні роки",
      semesters: "Семестри",
      departments: "Кафедри",
      rooms: "Аудиторії",
      specialties: "Спеціальності",
      cycleCommittees: "Циклові комісії",
      teachers: "Викладачі",
      disciplines: "Дисципліни",
      groups: "Групи",
      groupSemesters: "Групи по семестрах",
      bellSchedules: "Розклад дзвінків",
      /* Processing */
      studyPlans: "Навчальні плани",
      workloadDistributions: "Розподіл навантаження",
      workloadAssignments: "Призначення навантаження",
      /* Schedule */
      templateSettings: "Налаштування шаблону",
      restrictions: "Обмеження",
      templates: "Шаблони",
      generate: "Генерація",
      teacherSlotPreferences: "Переваги слотів",
      cycleCommitteeLabRooms: "Лаб. аудиторії",
    },

    generate: {
      title: "Генерація розкладу",
      semesterLabel: "Семестр",
      groupsLabel: "Групи",
      generateBtn: "Згенерувати розклад",
      generating: "Генерація…",
      success: "Розклад успішно згенеровано!",
      noSemester: "Будь ласка, оберіть семестр перед генерацією.",
      emptyHint: "← Оберіть параметри та натисніть «Згенерувати розклад»",
      saveBtn: "Зберегти шаблон",
      saveNamePlaceholder: "Назва шаблону (необов'язково)",
      savedSuccess: "Шаблон розкладу успішно збережено!",
      configTitle: "Налаштування генерації",
      settingsBtn: "Налаштування та обмеження",
      activateTitle: "Активувати розклад",
      activateConfirmText:
        "Розклад збережено. Перевірте параметри та активуйте його для використання.",
      activateSemester: "Семестр",
      activateGroups: "Груп охоплено",
      activateLessons: "Пар заплановано",
      activateBtn: "Активувати",
      activateSuccess: "Розклад успішно активовано!",
    },

    home: {
      activeTemplate: "Активний шаблон розкладу",
      noActiveTemplate:
        "Активний шаблон розкладу відсутній. Перейдіть до розділу «Розклад» для генерації.",
      activeYear: "Поточний навчальний рік",
      noActiveYear: "Активний навчальний рік відсутній.",
      createYear: "Створити навчальний рік",
    },

    roleTypes: {
      assistant: "Асистент",
      lector: "Лектор",
      lecture: "Лекція",
      practical: "Практичне",
      laboratory: "Лабораторне",
      exam: "Іспит",
      other: "Інше",
    },

    workload: {
      transfer: "До навантаження",
      transferGroup: "Оберіть групу",
      totalHours: "Усього годин",
      assignedHours: "Призначено",
      remainingHours: "Залишилось",
    },

    schedule: {
      numerator: "Чисельник",
      denominator: "Знаменник",
      viewTitle: "Перегляд розкладу",
      allGroups: "Всі групи",
      allTeachers: "Всі викладачі",
      groupFilter: "Група",
      teacherFilter: "Викладач",
      viewModeLabel: "Режим перегляду",
      groupView: "За групою",
      teacherView: "За викладачем",
      typeUnited: "Лекція",
      typeSplit: "Лаб",
      typeFlow: "Потік",
      typeLecture: "Лекція",
      typeLab: "Лаб",
      legendLecture: "Лекція",
      legendLab: "Лабораторна",
      legendFlow: "Потік",
      days: {
        monday: "Понеділок",
        tuesday: "Вівторок",
        wednesday: "Середа",
        thursday: "Четвер",
        friday: "П'ятниця",
      },
    },

    users: {
      email: "Електронна пошта",
      role: "Роль",
      firstName: "Ім'я",
      lastName: "Прізвище",
      password: "Пароль",
      activate: "Активувати",
      deactivate: "Деактивувати",
      register: "Зареєструвати",
      inviteLink: "Посилання для запрошення (поділіться з користувачем):",
      inviteLinkHint:
        "Поділіться цим посиланням з користувачем через месенджер або особисто. Воно дійсне 72 години.",
      noName: "(без імені)",
      pendingStatus: "Очікує активації",
      resetInvite: "Запросити",
      resetPassword: "Скинути пароль",
      deactivateConfirm:
        "Користувач {name} більше не зможе входити в систему. Підтвердити?",
    },
  },

  /* ── English ────────────────────────────────────────── */
  en: {
    nav: {
      home: "Home",
      reference: "Reference",
      processing: "Plans",
      schedule: "Schedule",
      academicSettings: "Academic Calendar",
      users: "Users",
      profile: "Profile",
    },

    auth: {
      loginTitle: "Sign in",
      loginSubtitle: "Enter your account credentials to continue.",
      emailLabel: "Email address",
      emailPlaceholder: "example@example.com",
      passwordLabel: "Password",
      loginBtn: "Sign in",
      logoutBtn: "Sign out",
      loginError: "Invalid email address or password.",
      accountDeactivated:
        "Your account has been deactivated. Please contact the administrator.",
      setPasswordTitle: "Set password",
      setPasswordBtn: "Set password",
      setPasswordOk: "Password set successfully. Signing in…",
      confirmPassword: "Confirm password",
      passwordMismatch: "Passwords do not match.",
      tokenInvalid: "The link is invalid or has expired.",
    },

    header: {
      academicYear: "Academic year",
      semester: "Semester",
      selectYear: "Select year",
      selectSemester: "Select semester",
      language: "Language",
    },

    lang: {
      ua: "Українська",
      en: "English",
    },

    sidebar: {
      mainNav: "Navigation",
      academicSettings: "Academic Calendar",
    },

    actions: {
      add: "Add",
      save: "Save",
      discard: "Discard",
      delete: "Delete",
      deleteSelected: "Delete selected",
      edit: "Edit",
      view: "View",
      close: "Close",
      cancel: "Cancel",
      confirm: "Confirm",
      generate: "Generate",
      activate: "Activate",
      deactivate: "Deactivate",
      optional: "optional",
      copy: "Copy",
      back: "Back",
      actions: "Actions",
    },

    messages: {
      saved: "Changes saved successfully.",
      saveError: "Failed to save data.",
      deleted: "Records deleted successfully.",
      deleteError: "Failed to delete records.",
      loadError: "Failed to load data.",
      confirmDelete: "Are you sure you want to delete the selected records?",
      unsavedChanges:
        "You have unsaved changes. Please save or discard them before leaving.",
      noAcademicYear:
        "No academic year selected. Please select or create an academic year.",
      noData: "No data",
      fillRequired: "Please fill in all required fields.",
      forbidden: "Access denied.",
      copied: "Copied!",
    },

    status: {
      active: "Active",
      inactive: "Inactive",
    },

    fields: {
      id: "ID",
      name: "Name",
      shortName: "Short name",
      email: "Email",
      role: "Role",
      password: "Password",
      startDate: "Start date",
      endDate: "End date",
      periodStart: "Period start",
      periodEnd: "Period end",
      isActive: "Active",
      isDeleted: "Deleted",
      createdAt: "Created at",
      modifiedAt: "Modified at",
      academicYear: "Academic year",
      semester: "Semester",
      department: "Department",
      specialty: "Specialty",
      cycleCommittee: "Cycle committee",
      teacher: "Teacher",
      discipline: "Discipline",
      group: "Group",
      lessonNumber: "Lesson number",
      startTime: "Start time",
      endTime: "End time",
      roomType: "Room type",
      specialtyCode: "Specialty code",
      isSplitting: "Split into subgroups",
      isSubvention: "Subvention",
      isContract: "Contract",
      educationStart: "Education start",
      educationEnd: "Education end",
      slotType: "Slot type",
      weekday: "Weekday",
      room: "Room",
      maxIdenticalLessonsPerDay: "Max identical lessons per day",
      maxStudyHoursPerDay: "Max study hours per day",
      maxTeacherHoursPerWeek: "Max teacher hours per week",
      maxGroupLessonHoursPerWeek: "Max group lesson hours per week",
      minGroupLessonsPerDay: "Min group lessons per day",
      maxGroupLessonsPerDay: "Max group lessons per day",
      maxTeacherLessonsPerDay: "Max teacher lessons per day",
      noGapsInGroupSchedule: "No gaps in group schedule",
      lessonsPerClass: "Lesson duration (acad. hrs)",
      studyDaysMask: "Study days mask",
      maxConsecutiveTeacherLessons: "Max consecutive teacher lessons",
      timePriority: "Time priority",
      allowFlowLessons: "Allow flow lessons",
      lecturesHours: "Lectures (hrs)",
      laboratoryHours: "Laboratory (hrs)",
      practicalHours: "Practical (hrs)",
      independentWork: "Independent work (hrs)",
      examHours: "Exam (hrs)",
      creditHours: "Credit (hrs)",
      classroomWork: "Classroom work (hrs)",
      semesterNumber: "Semester number",
      total: "Total",
      educationalComponent: "Educational component",
      studyPlan: "Study plan",
      workloadDistribution: "Workload distribution",
      assignedHours: "Assigned hours",
      roleType: "Role type",
      userId: "User ID",
      fullName: "Full name",
      firstName: "First name",
      lastName: "Last name",
      patronymic: "Patronymic",
    },

    entities: {
      academicYears: "Academic years",
      semesters: "Semesters",
      departments: "Departments",
      rooms: "Rooms",
      specialties: "Specialties",
      cycleCommittees: "Cycle committees",
      teachers: "Teachers",
      disciplines: "Disciplines",
      groups: "Groups",
      groupSemesters: "Group semesters",
      bellSchedules: "Bell schedules",
      studyPlans: "Study plans",
      workloadDistributions: "Workload distributions",
      workloadAssignments: "Workload assignments",
      scheduleTemplateSettings: "Template settings",
      scheduleRestrictions: "Schedule restrictions",
      scheduleTemplates: "Schedule templates",
      teacherSlotPreferences: "Teacher slot preferences",
      cycleCommitteeLabRooms: "Cycle committee lab rooms",
      users: "Users",
    },

    pages: {
      homeTitle: "Home",
      homeWelcome: "Welcome to EduPlanner",
      referenceTitle: "Reference",
      processingTitle: "Plans",
      scheduleTitle: "Schedule",
      academicSettingsTitle: "Academic Calendar",
      profileTitle: "Personal information",
      usersTitle: "User management",
      setPasswordTitle: "Set password",
    },

    roomTypes: {
      auditorium: "Auditorium",
      laboratory: "Laboratory",
    },

    lessonsPerClass: {
      one: "1 acad. hour",
      two: "2 acad. hours",
      three: "3 acad. hours",
    },

    timePriority: {
      none: "No priority",
      morning: "First half of day",
      afternoon: "Second half of day",
    },

    weekdays: {
      monday: "Monday",
      tuesday: "Tuesday",
      wednesday: "Wednesday",
      thursday: "Thursday",
      friday: "Friday",
      saturday: "Saturday",
    },

    slotTypes: {
      preferred: "Preferred",
      forbidden: "Forbidden",
    },

    roles: {
      admin: "Administrator",
      dean: "Dean office",
      guest: "Guest",
      teacher: "Teacher",
    },

    units: {
      lessons: "lessons",
      hours: "hrs",
      consecutive: "in a row",
    },

    weekdaysShort: {
      monday: "Mo",
      tuesday: "Tu",
      wednesday: "We",
      thursday: "Th",
      friday: "Fr",
      saturday: "Sa",
      sunday: "Su",
    },

    modal: {
      unsaved: "Unsaved changes",
      saved: "All changes saved",

      tabFormat: "Lesson format",
      tabCalendar: "Calendar",
      tabDaily: "Daily limits",
      tabWeekly: "Weekly limits",
      tabLoad: "Load",
      tabTime: "Time & calendar",
      tabSpace: "Rooms",
      tabOrg: "Organizational",

      lessonDuration: "Lesson duration",
      lessonDurationDesc: "Number of academic hours in one lesson.",
      studyDays: "Study days",
      studyDaysDesc: "Days of the week when the schedule is generated.",
      maxIdentical: "Max identical lessons per day",
      maxIdenticalDesc:
        "How many lessons of the same subject can be on one day.",
      maxHoursDay: "Max study hours per day",
      maxHoursDayDesc:
        "Total upper limit of classroom hours in a group's schedule per day.",
      maxTeacherWeek: "Max teacher hours per week",
      maxTeacherWeekDesc:
        "Upper limit of classroom workload for one teacher per week.",
      maxGroupWeek: "Max group lesson hours per week",
      maxGroupWeekDesc:
        "Total upper limit of hours in one group's schedule per week.",
      maxGroupDay: "Max lessons per day for students",
      maxGroupDayDesc:
        "Upper limit on lessons in a group's schedule per study day.",
      maxTeacherDay: "Max lessons per day for teachers",
      maxTeacherDayDesc:
        "Maximum number of lessons a teacher can conduct per day.",
      consecutive: "Consecutive teacher work",
      consecutiveDesc:
        "Max consecutive lessons without a break for one teacher (0 = no limit).",
      minGroupDay: "Min group lessons per day",
      minGroupDayDesc:
        "If a group has lessons on this day — minimum number (0 = no limit).",
      gapPrevention: "Gap prevention",
      gapPreventionDesc:
        "Minimize empty hours between lessons for students and teachers.",
      gapOff: "Disabled",
      gapOn: "Fully exclude",
      forbiddenHours: "Forbidden teacher hours",
      forbiddenHoursDesc: "Times when a teacher cannot conduct lessons.",
      teacherPrefs: "Teacher preferences",
      teacherPrefsDesc:
        "Priority hours when a teacher wants to conduct lessons.",
      timePriority: "Lesson time priority",
      timePriorityDesc:
        "Preferred half of the day for placing lessons in the schedule.",
      labRooms: "Lab rooms for disciplines",
      labRoomsDesc:
        "Select computer labs where laboratory work can be conducted.",
      flowLessons: "Flow lessons (shared lectures)",
      flowLessonsDesc:
        "Allow one teacher to conduct a lesson simultaneously for multiple groups.",
      flowOn: "Allowed",
      flowOff: "Forbidden",
      configure: "Configure",
      back: "Back",
      noSettings: "No settings",
      selectTeacher: "Select a teacher",
      lessonNum: "Lesson #",
      chooseCycleCommittee: "Choose cycle committee…",
      forbiddenTitle: "Forbidden teacher hours",
      preferredTitle: "Teacher preferences",
    },

    setupModal: {
      title: "Academic year setup",
      description:
        "No academic year found. Please fill in the details to create a new academic year and first semester.",
      yearName: "Academic year name",
      yearStart: "Academic year start",
      yearEnd: "Academic year end",
      semesterStart: "Semester start",
      semesterEnd: "Semester end",
      createBtn: "Create",
    },

    tabs: {
      /* Reference */
      academicYears: "Academic years",
      semesters: "Semesters",
      departments: "Departments",
      rooms: "Rooms",
      specialties: "Specialties",
      cycleCommittees: "Cycle committees",
      teachers: "Teachers",
      disciplines: "Disciplines",
      groups: "Groups",
      groupSemesters: "Group semesters",
      bellSchedules: "Bell schedules",
      /* Processing */
      studyPlans: "Study plans",
      workloadDistributions: "Workload distributions",
      workloadAssignments: "Workload assignments",
      /* Schedule */
      templateSettings: "Template settings",
      restrictions: "Restrictions",
      templates: "Templates",
      generate: "Generate",
      teacherSlotPreferences: "Slot preferences",
      cycleCommitteeLabRooms: "Lab rooms",
    },

    generate: {
      title: "Schedule generation",
      semesterLabel: "Semester",
      groupsLabel: "Groups",
      generateBtn: "Generate schedule",
      generating: "Generating…",
      success: "Schedule generated successfully!",
      noSemester: "Please select a semester before generating.",
      emptyHint: "← Select parameters and click «Generate schedule»",
      saveBtn: "Save template",
      saveNamePlaceholder: "Template name (optional)",
      savedSuccess: "Schedule template saved successfully!",
      configTitle: "Generation settings",
      settingsBtn: "Settings & restrictions",
      activateTitle: "Activate schedule",
      activateConfirmText:
        "The schedule has been saved. Review the parameters and activate it for use.",
      activateSemester: "Semester",
      activateGroups: "Groups covered",
      activateLessons: "Lessons planned",
      activateBtn: "Activate",
      activateSuccess: "Schedule activated successfully!",
    },

    home: {
      activeTemplate: "Active schedule template",
      noActiveTemplate:
        "No active schedule template. Go to the «Schedule» section to generate one.",
      activeYear: "Current academic year",
      noActiveYear: "No active academic year.",
      createYear: "Create academic year",
    },

    roleTypes: {
      assistant: "Assistant",
      lector: "Lecturer",
      lecture: "Lecture",
      practical: "Practical",
      laboratory: "Laboratory",
      exam: "Exam",
      other: "Other",
    },

    workload: {
      transfer: "To workload",
      transferGroup: "Select group",
      totalHours: "Total hours",
      assignedHours: "Assigned",
      remainingHours: "Remaining",
    },

    schedule: {
      numerator: "Numerator",
      denominator: "Denominator",
      viewTitle: "Schedule view",
      allGroups: "All groups",
      allTeachers: "All teachers",
      groupFilter: "Group",
      teacherFilter: "Teacher",
      viewModeLabel: "View mode",
      groupView: "By group",
      teacherView: "By teacher",
      typeUnited: "Lecture",
      typeSplit: "Lab",
      typeFlow: "Flow",
      typeLecture: "Lecture",
      typeLab: "Lab",
      legendLecture: "Lecture",
      legendLab: "Laboratory",
      legendFlow: "Flow",
      days: {
        monday: "Monday",
        tuesday: "Tuesday",
        wednesday: "Wednesday",
        thursday: "Thursday",
        friday: "Friday",
      },
    },

    users: {
      email: "Email",
      role: "Role",
      firstName: "First name",
      lastName: "Last name",
      password: "Password",
      activate: "Activate",
      deactivate: "Deactivate",
      register: "Register",
      inviteLink: "Invite link (share with user):",
      inviteLinkHint:
        "Share this link with the user via messenger or in person. It expires in 72 hours.",
      noName: "(no name)",
      pendingStatus: "Pending",
      resetInvite: "Invite",
      resetPassword: "Reset password",
      deactivateConfirm:
        "User {name} will no longer be able to log in. Confirm?",
    },
  },
};

/* =========================================================
   Runtime language state
   ========================================================= */

/** @type {'ua'|'en'} */
let currentLang = localStorage.getItem("lang") || "ua";

/**
 * Resolve a dot-notation key against a translation object.
 * @param {object} obj
 * @param {string} key  e.g. "actions.save"
 * @returns {string|undefined}
 */
function _resolve(obj, key) {
  return key.split(".").reduce(function (acc, part) {
    return acc != null && typeof acc === "object" ? acc[part] : undefined;
  }, obj);
}

/**
 * Get a translated string by dot-notation key.
 * Falls back to the English translation if the current language lacks the key.
 * Falls back to the key itself if neither language has it.
 *
 * @param {string} key  e.g. "actions.save"
 * @returns {string}
 */
function t(key) {
  var value = _resolve(TRANSLATIONS[currentLang], key);
  if (value == null || typeof value !== "string") {
    value = _resolve(TRANSLATIONS["en"], key);
  }
  if (value == null || typeof value !== "string") {
    // Return last segment of key as readable fallback
    var parts = key.split(".");
    return parts[parts.length - 1];
  }
  return value;
}

/**
 * Switch the active language and refresh all translated DOM elements.
 * @param {'ua'|'en'} lang
 */
function setLang(lang) {
  if (!TRANSLATIONS[lang]) {
    console.warn("[i18n] Unknown language:", lang);
    return;
  }
  currentLang = lang;
  localStorage.setItem("lang", lang);
  applyTranslations();
}

/**
 * Update every element that carries a [data-i18n] attribute.
 *
 * Supported attribute forms:
 *   data-i18n="key"              → sets element.textContent
 *   data-i18n-placeholder="key" → sets element.placeholder
 *   data-i18n-title="key"       → sets element.title
 *   data-i18n-value="key"       → sets element.value (for buttons)
 *   data-i18n-aria="key"        → sets element.ariaLabel
 */
function applyTranslations() {
  // textContent
  document.querySelectorAll("[data-i18n]").forEach(function (el) {
    var key = el.getAttribute("data-i18n");
    var text = t(key);
    // Avoid touching elements that contain child elements (only text nodes)
    if (el.children.length === 0) {
      el.textContent = text;
    } else {
      // Update the first text node only
      for (var i = 0; i < el.childNodes.length; i++) {
        if (el.childNodes[i].nodeType === Node.TEXT_NODE) {
          el.childNodes[i].textContent = text;
          break;
        }
      }
    }
  });

  // placeholder
  document.querySelectorAll("[data-i18n-placeholder]").forEach(function (el) {
    el.placeholder = t(el.getAttribute("data-i18n-placeholder"));
  });

  // title (tooltip)
  document.querySelectorAll("[data-i18n-title]").forEach(function (el) {
    el.title = t(el.getAttribute("data-i18n-title"));
  });

  // value (e.g. submit buttons)
  document.querySelectorAll("[data-i18n-value]").forEach(function (el) {
    el.value = t(el.getAttribute("data-i18n-value"));
  });

  // aria-label
  document.querySelectorAll("[data-i18n-aria]").forEach(function (el) {
    el.setAttribute("aria-label", t(el.getAttribute("data-i18n-aria")));
  });

  // Sync sidebar lang buttons
  document.querySelectorAll(".lang-btn[data-lang]").forEach(function (btn) {
    btn.classList.toggle(
      "active",
      btn.getAttribute("data-lang") === currentLang,
    );
  });

  // Notify the rest of the app (optional listener pattern)
  window.dispatchEvent(
    new CustomEvent("langchange", { detail: { lang: currentLang } }),
  );
}

/* =========================================================
   Auto-init
   ========================================================= */

// Apply translations once the DOM is ready.
if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", applyTranslations);
} else {
  applyTranslations();
}
