<template>
    <div v-if="modelValue" class="modal-overlay" @click="closeOnOverlay">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <slot name="title"></slot>
        </div>
        <div class="modal-body">
          <slot></slot>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  const props = defineProps({
    modelValue: Boolean,
    closeOnOutsideClick: {
      type: Boolean,
      default: true
    }
  })
  
  const emit = defineEmits(['update:modelValue'])
  
  const closeOnOverlay = () => {
    if (props.closeOnOutsideClick) {
      emit('update:modelValue', false)
    }
  }
  </script>
  
  <style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 100;
  }
  
  .modal-content {
    background: white;
    border-radius: 8px;
    padding: 24px;
    width: 90%;
    max-width: 500px;
    max-height: 90vh;
    overflow-y: auto;
  }
  
  .modal-header {
    margin-bottom: 16px;
  }
  
  .modal-header h3 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
  }
  </style>