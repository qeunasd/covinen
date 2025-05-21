function stripEmptyInputs(form) {
  const inputs = form.querySelectorAll("input[name], select[name], textarea[name]");
  inputs.forEach(input => {
    const isCheckboxOrRadio = input.type === "checkbox" || input.type === "radio";
    if (isCheckboxOrRadio && !input.checked) {
      input.removeAttribute("name");
    } else if (!isCheckboxOrRadio && !input.value.trim()) {
      input.removeAttribute("name");
    }
  });
}

function toggleRowHighlight(checkbox) {
  const row = checkbox.closest("tr");

  if (checkbox.checked) {
    row.classList.add("bg-yellow-100", "ring-1", "ring-yellow-400", "hover:bg-yellow-200");
    row.classList.remove("hover:bg-gray-50");
  } else {
    row.classList.remove("bg-yellow-100", "ring-1", "ring-yellow-400", "hover:bg-yellow-200");
    row.classList.add("hover:bg-gray-50");
  }
}

function toggleAll(source) {
  const checkboxes = document.querySelectorAll('input[name="ids"]');
  checkboxes.forEach(cb => {
    cb.checked = source.checked;
    toggleRowHighlight(cb);
  });
}