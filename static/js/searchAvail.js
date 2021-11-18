function SearchAvailJson(roomId, token) {
    document.getElementById("check-availability-button").addEventListener('click', function () {
        // notify("This is a test message!","success");
        // notifyModal("Success Text","<em>Hello, World!</em>",'success',"My Text For The Button")
        // attention.toast({msg:"Hello, world!",icon:'error'})
        // attention.success({msg:"Hello, world!"});
        //attention.error({msg:"Hello, world!"});
        let html = `
            <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">

                  <div class="row container-fluid" id="reservation-date-modal">
                        <div class="col">
                        <label for="start_date"> Starting Date: </label>
                        <input type="text" class="form-control mt-2" autocomplete="off" id="start_date" name="start_date" placeholder="Arrival" disabled required>
                        </div>
                        <div class="col">
                        <label for="end_date">Ending Date: </label>
                        <input type="text" class="form-control mt-2" autocomplete="off" id="end_date" name="end_date" placeholder="Departure" disabled required>
                        </div>
                  </div>


          </form>
          `;
        attention.custom({
            msg: html,
            title: "Form Test!",
            willOpen() {
                const elem = document.getElementById('reservation-date-modal');
                const rp = new DateRangePicker(elem, {
                    format: 'yyyy-mm-dd',
                    showOnFocus: true,
                    minDate: new Date(),
                })
            },
            preConfirm: () => {
                return [
                    document.getElementById('start_date').value,
                    document.getElementById('end_date').value
                ]
            },
            didOpen: () => {
                document.getElementById('start_date').removeAttribute('disabled');
                document.getElementById('end_date').removeAttribute('disabled');
            },
            callback: function (result) {
                console.log("Called");
                let form = document.getElementById("check-availability-form");
                let formData = new FormData(form);
                formData.append("csrf_token", token)
                formData.append("room_id", roomId)
                //    here we can write some javascript code that will perform an ajax request
                fetch('/search-availability-json', {
                    method: 'post',
                    body: formData
                }).then(response => response.json())
                    .then(data => {
                        if (data.ok) {
                            attention.custom({
                                icon: 'success',
                                msg: '<p>Room is available</p>' +
                                    '<p><a href="/book-room?id=' + data.room_id
                                    + '&s=' + data.start_date
                                    + '&e=' + data.end_date
                                    + '" id ="submit-reservation-button"'
                                    + ' class = "btn btn-primary" >' + 'Book Now!' + '</a></p >',
                                showConfirmButton: false,
                            })
                        } else {
                            attention.error(
                                {
                                    msg: "No Availability!",
                                }
                            )
                        }
                    })
            }
        });
    })
}