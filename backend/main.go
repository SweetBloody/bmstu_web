package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	authHandler "github.com/SweetBloody/bmstu_web/backend/internal/pkg/auth/delivery/http"
	driverHandler "github.com/SweetBloody/bmstu_web/backend/internal/pkg/driver/delivery/http"
	grandPrixHandler "github.com/SweetBloody/bmstu_web/backend/internal/pkg/grand_prix/delivery/http"
	qualHandler "github.com/SweetBloody/bmstu_web/backend/internal/pkg/qual_result/delivery/http"
	raceHandler "github.com/SweetBloody/bmstu_web/backend/internal/pkg/race_result/delivery/http"
	teamHandler "github.com/SweetBloody/bmstu_web/backend/internal/pkg/team/delivery/http"
	trackHandler "github.com/SweetBloody/bmstu_web/backend/internal/pkg/track/delivery/http"
	userHandler "github.com/SweetBloody/bmstu_web/backend/internal/pkg/user/delivery/http"

	driverRepository "github.com/SweetBloody/bmstu_web/backend/internal/pkg/driver/repository/postgresql"
	grandPrixRepository "github.com/SweetBloody/bmstu_web/backend/internal/pkg/grand_prix/repository/postgresql"
	qualRepository "github.com/SweetBloody/bmstu_web/backend/internal/pkg/qual_result/repository/postgresql"
	raceRepository "github.com/SweetBloody/bmstu_web/backend/internal/pkg/race_result/repository/postgresql"
	teamRepository "github.com/SweetBloody/bmstu_web/backend/internal/pkg/team/repository/postgresql"
	trackRepository "github.com/SweetBloody/bmstu_web/backend/internal/pkg/track/repository/postgresql"
	userRepository "github.com/SweetBloody/bmstu_web/backend/internal/pkg/user/repository/postgresql"

	driverUsecase "github.com/SweetBloody/bmstu_web/backend/internal/pkg/driver/usecase"
	grandPrixUsecase "github.com/SweetBloody/bmstu_web/backend/internal/pkg/grand_prix/usecase"
	qualUsecase "github.com/SweetBloody/bmstu_web/backend/internal/pkg/qual_result/usecase"
	raceUsecase "github.com/SweetBloody/bmstu_web/backend/internal/pkg/race_result/usecase"
	teamUsecase "github.com/SweetBloody/bmstu_web/backend/internal/pkg/team/usecase"
	trackUsecase "github.com/SweetBloody/bmstu_web/backend/internal/pkg/track/usecase"
	userUsecase "github.com/SweetBloody/bmstu_web/backend/internal/pkg/user/usecase"

	_ "github.com/SweetBloody/bmstu_web/backend/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/SweetBloody/bmstu_web/backend/internal/app/middleware"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// @title FormulOne Web-Server
// @version 1.0
// @description API Server for F1 Grand-Prix info
// @termsOfService  http://swagger.io/terms/

// @host localhost:5259
// @BasePath /
func main() {
	params := fmt.Sprintf("user=default_admin dbname=back password=12345678 host=%s port=%s sslmode=disable", os.Getenv("PG_HOST"), os.Getenv("PG_PORT"))
	db, err := sqlx.Connect("postgres", params)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driverRepo := driverRepository.NewPsqlDriverRepository(db)
	teamRepo := teamRepository.NewPsqlTeamRepository(db)
	trackRepo := trackRepository.NewPsqlTrackRepository(db)
	gpRepo := grandPrixRepository.NewPsqlGPRepository(db)
	raceRepo := raceRepository.NewPsqlRaceResultRepository(db)
	qualRepo := qualRepository.NewPsqlQualResultRepository(db)
	userRepo := userRepository.NewPsqlUserRepository(db)

	driverUcase := driverUsecase.NewDriverUsecase(driverRepo)
	teamUcase := teamUsecase.NewTeamUsecase(teamRepo)
	trackUcase := trackUsecase.NewTrackUsecase(trackRepo)
	gpUcase := grandPrixUsecase.NewGrandPrixUsecase(gpRepo)
	raceUcase := raceUsecase.NewRaceResultUsecase(raceRepo)
	qualUcase := qualUsecase.NewQualResultUsecase(qualRepo)
	userUcase := userUsecase.NewUserUsecase(userRepo)

	m := mux.NewRouter()

	driverHandler.NewDriverHandler(m, driverUcase, raceUcase)
	teamHandler.NewTeamHandler(m, teamUcase)
	trackHandler.NewTrackHandler(m, trackUcase)
	grandPrixHandler.NewDriverHandler(m, gpUcase, raceUcase, qualUcase)
	raceHandler.NewRaceResultHandler(m, raceUcase)
	qualHandler.NewQualResultHandler(m, qualUcase)
	authHandler.NewAuthHandler(m, userUcase)
	userHandler.NewUserHandler(m, userUcase)

	m.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)
	mMiddleware := middleware.LogMiddleware(m)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", handlers.CORS()(mMiddleware))
}
