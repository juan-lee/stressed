/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

const defaultJobFile = `# stress-ng jobfile
    run sequential   # run stressors sequentially
    verbose          # verbose output
    metrics-brief    # show metrics at end of run
    timeout 60s      # stop each stressor after 60 seconds

    ################################################################################
    # runs 4 cpu stressors, 2 io stressors and 1 vm stressor using 1GB of virtual
    # memory.
    ################################################################################
    cpu 4
    io 2
    vm 1
    vm-bytes 1G

    ################################################################################
    # run 8 virtual memory stressors that combined use 80% of the available memory.
    # Thus each stressor uses 10% of the available memory.
    ################################################################################
    #vm 8
    #vm-bytes 80%

    ################################################################################
    # runs 2 instances of the mixed I/O stressors using a total of 10% of the 
    # available file system space. Each stressor will use 5% of the available file 
    # system space.
    ################################################################################
    # iomix 2
    # iomix-bytes 10%
`
